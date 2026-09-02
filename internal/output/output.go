// Package output renders results in the format the user asked for.
//
// text is the default and is meant to be read; json and yaml are meant to be
// piped. The text renderer works by reflection over whatever the SDK returned,
// so a new operation needs no rendering code of its own.
package output

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"

	"gopkg.in/yaml.v3"
)

// Format is a rendering style.
type Format string

const (
	Text Format = "text"
	JSON Format = "json"
	YAML Format = "yaml"
)

// ParseFormat validates a --output value.
func ParseFormat(s string) (Format, error) {
	switch Format(strings.ToLower(strings.TrimSpace(s))) {
	case "", Text:
		return Text, nil
	case JSON:
		return JSON, nil
	case YAML:
		return YAML, nil
	}
	return "", fmt.Errorf("unknown output format %q: use text, json or yaml", s)
}

// Printer writes results.
type Printer struct {
	Out    io.Writer
	Err    io.Writer
	Format Format
	// NoHeaders suppresses the header row in text tables, for piping into
	// awk and friends.
	NoHeaders bool
}

// Value renders a single result.
func (p *Printer) Value(v any) error {
	if v == nil || isNilPointer(v) {
		return nil
	}
	switch p.Format {
	case JSON:
		return p.writeJSON(v)
	case YAML:
		return p.writeYAML(v)
	}
	return p.writeFields(v)
}

// Page renders one page of a list, then a hint when more remain.
func (p *Printer) Page(page any) error {
	items, meta := pageParts(page)
	if !items.IsValid() {
		return p.Value(page)
	}
	switch p.Format {
	case JSON:
		return p.writeJSON(items.Interface())
	case YAML:
		return p.writeYAML(items.Interface())
	}
	if err := p.writeTable(items); err != nil {
		return err
	}
	// Only worth saying in text mode, and only when it changes what to do
	// next.
	if meta.hasMore && p.Err != nil {
		fmt.Fprintf(p.Err, "\nMore results available. Re-run with --all, or --marker %s\n", meta.marker)
	}
	return nil
}

// Iter renders every item an iterator yields, streaming in text mode so a
// long walk shows progress rather than buffering to the end.
func (p *Printer) Iter(seq any) error {
	v := reflect.ValueOf(seq)
	if v.Kind() != reflect.Func {
		return fmt.Errorf("output: not an iterator: %T", seq)
	}
	elem := v.Type().In(0).In(0)
	collected := reflect.MakeSlice(reflect.SliceOf(elem), 0, 0)

	var iterErr error
	yieldType := v.Type().In(0)
	yield := reflect.MakeFunc(yieldType, func(args []reflect.Value) []reflect.Value {
		if len(args) > 1 && !args[1].IsNil() {
			iterErr = args[1].Interface().(error)
			return []reflect.Value{reflect.ValueOf(false)}
		}
		collected = reflect.Append(collected, args[0])
		return []reflect.Value{reflect.ValueOf(true)}
	})
	v.Call([]reflect.Value{yield})
	if iterErr != nil {
		return iterErr
	}

	switch p.Format {
	case JSON:
		return p.writeJSON(collected.Interface())
	case YAML:
		return p.writeYAML(collected.Interface())
	}
	return p.writeTable(collected)
}

// Stream copies raw bytes through, for an object body or a screenshot.
func (p *Printer) Stream(r io.ReadCloser) error {
	defer r.Close()
	_, err := io.Copy(p.Out, r)
	return err
}

// Done reports a successful action that returned nothing.
func (p *Printer) Done(msg string) {
	if p.Format == Text && msg != "" {
		fmt.Fprintln(p.Out, msg)
	}
}

func (p *Printer) writeJSON(v any) error {
	enc := json.NewEncoder(p.Out)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func (p *Printer) writeYAML(v any) error {
	// Round-trip through JSON so the output uses the API's field names rather
	// than Go's, which is what someone comparing against the docs expects.
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	var generic any
	if err := json.Unmarshal(data, &generic); err != nil {
		return err
	}
	enc := yaml.NewEncoder(p.Out)
	enc.SetIndent(2)
	if err := enc.Encode(generic); err != nil {
		return err
	}
	return enc.Close()
}

// writeFields prints one object as aligned key/value lines.
func (p *Printer) writeFields(v any) error {
	m, err := toMap(v)
	if err != nil {
		return err
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	// Identity first: it is what a reader looks for.
	sort.SliceStable(keys, func(i, j int) bool {
		return fieldRank(keys[i]) < fieldRank(keys[j])
	})

	tw := tabwriter.NewWriter(p.Out, 0, 0, 2, ' ', 0)
	for _, k := range keys {
		fmt.Fprintf(tw, "%s:\t%s\n", k, scalarString(m[k]))
	}
	return tw.Flush()
}

// writeTable prints a slice of objects as a table.
func (p *Printer) writeTable(items reflect.Value) error {
	if items.Len() == 0 {
		if p.Err != nil {
			fmt.Fprintln(p.Err, "No results.")
		}
		return nil
	}
	rows := make([]map[string]any, 0, items.Len())
	for i := 0; i < items.Len(); i++ {
		m, err := toMap(items.Index(i).Interface())
		if err != nil {
			return err
		}
		rows = append(rows, m)
	}
	cols := chooseColumns(rows)

	tw := tabwriter.NewWriter(p.Out, 0, 0, 3, ' ', 0)
	if !p.NoHeaders {
		header := make([]string, len(cols))
		for i, c := range cols {
			header[i] = strings.ToUpper(c)
		}
		fmt.Fprintln(tw, strings.Join(header, "\t"))
	}
	for _, row := range rows {
		cells := make([]string, len(cols))
		for i, c := range cols {
			cells[i] = scalarString(row[c])
		}
		fmt.Fprintln(tw, strings.Join(cells, "\t"))
	}
	return tw.Flush()
}

// preferredColumns are shown first when a result has them. The order is what
// someone scanning a list looks for: which one is this, what is it called,
// what state is it in.
var preferredColumns = []string{
	"id", "name", "status", "state", "vm_state", "power_state",
	"type", "region", "cidr", "address", "created_at",
}

func fieldRank(name string) int {
	for i, p := range preferredColumns {
		if name == p {
			return i
		}
	}
	return len(preferredColumns) + 1
}

// chooseColumns picks a readable subset. A table with thirty columns is not a
// table, so this shows the identifying fields and fills up to a cap.
func chooseColumns(rows []map[string]any) []string {
	const maxColumns = 7

	present := map[string]bool{}
	for _, row := range rows {
		for k, v := range row {
			if isScalar(v) {
				present[k] = true
			}
		}
	}
	var cols []string
	for _, c := range preferredColumns {
		if present[c] {
			cols = append(cols, c)
			delete(present, c)
		}
	}
	// A few fields are real but rarely what someone scanning a table wants,
	// and they are wide enough to push out a field that is. They go last, so
	// they appear only when there is room.
	deprioritised := map[string]bool{
		"crn": true, "arn": true, "updated_at": true, "organization_id": true,
		"account_id": true, "region_id": true, "project_id": true,
	}
	var rest, tail []string
	for k := range present {
		if deprioritised[k] {
			tail = append(tail, k)
			continue
		}
		rest = append(rest, k)
	}
	sort.Strings(rest)
	sort.Strings(tail)
	rest = append(rest, tail...)
	for _, k := range rest {
		if len(cols) >= maxColumns {
			break
		}
		cols = append(cols, k)
	}
	if len(cols) == 0 {
		cols = []string{"value"}
	}
	return cols
}

// toMap renders a value through its JSON form, so the keys are the API's
// field names rather than Go's.
func toMap(v any) (map[string]any, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		// Not an object: a bare string or number is still worth printing.
		var scalar any
		if err := json.Unmarshal(data, &scalar); err != nil {
			return nil, err
		}
		return map[string]any{"value": scalar}, nil
	}
	return m, nil
}

func isScalar(v any) bool {
	switch v.(type) {
	case string, float64, bool, nil:
		return true
	}
	return false
}

func scalarString(v any) string {
	switch t := v.(type) {
	case nil:
		return ""
	case string:
		return t
	case bool:
		return strconv.FormatBool(t)
	case float64:
		if t == float64(int64(t)) {
			return strconv.FormatInt(int64(t), 10)
		}
		return strconv.FormatFloat(t, 'f', -1, 64)
	case []any:
		parts := make([]string, len(t))
		for i, e := range t {
			parts[i] = scalarString(e)
		}
		return strings.Join(parts, ",")
	}
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprint(v)
	}
	return string(data)
}

type pageMeta struct {
	hasMore bool
	marker  string
}

// pageParts reaches into a *basaltic.Page[T] without naming its element type.
func pageParts(page any) (reflect.Value, pageMeta) {
	v := reflect.ValueOf(page)
	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return reflect.Value{}, pageMeta{}
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return reflect.Value{}, pageMeta{}
	}
	items := v.FieldByName("Items")
	if !items.IsValid() {
		return reflect.Value{}, pageMeta{}
	}
	meta := pageMeta{}
	if f := v.FieldByName("HasMore"); f.IsValid() && f.Kind() == reflect.Bool {
		meta.hasMore = f.Bool()
	}
	if f := v.FieldByName("Marker"); f.IsValid() && f.Kind() == reflect.String {
		meta.marker = f.String()
	}
	return items, meta
}

func isNilPointer(v any) bool {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Pointer, reflect.Map, reflect.Slice, reflect.Interface:
		return rv.IsNil()
	}
	return false
}
