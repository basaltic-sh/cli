package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

type widget struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Status    string            `json:"status"`
	CRN       string            `json:"crn"`
	Tags      map[string]string `json:"tags,omitempty"`
	SizeGB    int               `json:"size_gb"`
	Encrypted bool              `json:"encrypted"`
}

func newPrinter(f Format) (*Printer, *bytes.Buffer, *bytes.Buffer) {
	out, errBuf := &bytes.Buffer{}, &bytes.Buffer{}
	return &Printer{Out: out, Err: errBuf, Format: f}, out, errBuf
}

func TestTableShowsIdentifyingColumnsFirst(t *testing.T) {
	p, out, _ := newPrinter(Text)
	items := []widget{
		{ID: "w-1", Name: "alpha", Status: "active", CRN: "crn:x:::widget/w-1", SizeGB: 10},
		{ID: "w-2", Name: "beta", Status: "pending", CRN: "crn:x:::widget/w-2", SizeGB: 20},
	}
	if err := p.Iter(sliceSeq(items)); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	header := lines[0]

	// Identity leads, because it is what someone scanning a list looks for.
	if !strings.HasPrefix(header, "ID") {
		t.Errorf("header does not start with ID: %q", header)
	}
	idx := func(s string) int { return strings.Index(header, s) }
	if !(idx("ID") < idx("NAME") && idx("NAME") < idx("STATUS")) {
		t.Errorf("columns are out of order: %q", header)
	}
	// CRN is derivable from the id and wide enough to push out something
	// useful, so it goes last.
	if idx("CRN") < idx("SIZE_GB") {
		t.Errorf("CRN should come after the ordinary fields: %q", header)
	}
	if len(lines) != 3 {
		t.Errorf("got %d lines, want a header and two rows", len(lines))
	}
}

func TestNoHeaders(t *testing.T) {
	p, out, _ := newPrinter(Text)
	p.NoHeaders = true
	if err := p.Iter(sliceSeq([]widget{{ID: "w-1", Name: "alpha"}})); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out.String(), "ID") {
		t.Errorf("header printed despite NoHeaders: %q", out.String())
	}
}

func TestEmptyResultSaysSoOnStderr(t *testing.T) {
	p, out, errBuf := newPrinter(Text)
	if err := p.Iter(sliceSeq([]widget{})); err != nil {
		t.Fatal(err)
	}
	if out.Len() != 0 {
		t.Errorf("stdout should stay empty so a pipe sees nothing: %q", out.String())
	}
	if !strings.Contains(errBuf.String(), "No results") {
		t.Errorf("nothing said about an empty result: %q", errBuf.String())
	}
}

func TestJSONUsesTheAPIFieldNames(t *testing.T) {
	p, out, _ := newPrinter(JSON)
	if err := p.Value(&widget{ID: "w-1", Name: "alpha", SizeGB: 10}); err != nil {
		t.Fatal(err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(out.Bytes(), &decoded); err != nil {
		t.Fatalf("output is not JSON: %v", err)
	}
	if _, ok := decoded["size_gb"]; !ok {
		t.Errorf("JSON uses Go field names rather than the API's: %v", decoded)
	}
}

func TestYAMLUsesTheAPIFieldNames(t *testing.T) {
	p, out, _ := newPrinter(YAML)
	if err := p.Value(&widget{ID: "w-1", SizeGB: 10}); err != nil {
		t.Fatal(err)
	}
	// Round-tripping through JSON is what makes this true; yaml.v3 would
	// otherwise lowercase the Go names and produce "sizegb".
	if !strings.Contains(out.String(), "size_gb:") {
		t.Errorf("YAML does not use the API's field names: %q", out.String())
	}
}

func TestPageReportsMoreResultsOnStderr(t *testing.T) {
	p, out, errBuf := newPrinter(Text)
	page := struct {
		Items   []widget
		HasMore bool
		Marker  string
	}{Items: []widget{{ID: "w-1"}}, HasMore: true, Marker: "next-page"}

	if err := p.Page(&page); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "w-1") {
		t.Errorf("the page's items were not printed: %q", out.String())
	}
	// The hint goes to stderr so it cannot corrupt a pipe.
	if !strings.Contains(errBuf.String(), "--all") || !strings.Contains(errBuf.String(), "next-page") {
		t.Errorf("no usable hint about the remaining pages: %q", errBuf.String())
	}
	if strings.Contains(out.String(), "--all") {
		t.Errorf("the hint reached stdout: %q", out.String())
	}
}

func TestPageInJSONPrintsOnlyTheItems(t *testing.T) {
	p, out, _ := newPrinter(JSON)
	page := struct {
		Items   []widget
		HasMore bool
		Marker  string
	}{Items: []widget{{ID: "w-1"}}, HasMore: true, Marker: "m"}

	if err := p.Page(&page); err != nil {
		t.Fatal(err)
	}
	// A consumer piping to jq wants the array, not an envelope around it.
	var decoded []map[string]any
	if err := json.Unmarshal(out.Bytes(), &decoded); err != nil {
		t.Fatalf("output is not a JSON array: %v (%s)", err, out.String())
	}
	if len(decoded) != 1 || decoded[0]["id"] != "w-1" {
		t.Errorf("unexpected items: %v", decoded)
	}
}

func TestParseFormat(t *testing.T) {
	for in, want := range map[string]Format{"": Text, "text": Text, "JSON": JSON, "yaml": YAML} {
		got, err := ParseFormat(in)
		if err != nil || got != want {
			t.Errorf("ParseFormat(%q) = (%q, %v), want %q", in, got, err, want)
		}
	}
	if _, err := ParseFormat("xml"); err == nil {
		t.Error("ParseFormat accepted an unknown format")
	}
}

func sliceSeq[T any](items []T) func(func(T, error) bool) {
	return func(yield func(T, error) bool) {
		for _, it := range items {
			if !yield(it, nil) {
				return
			}
		}
	}
}
