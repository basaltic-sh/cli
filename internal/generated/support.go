// Package generated holds the command tree built from the SDK manifest.
//
// Every file ending in _gen.go is generated and must not be edited; this file
// is the hand-written support the generated code calls.
package generated

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/basaltic-sh/cli/internal/progress"

	"gopkg.in/yaml.v3"
)

// loadBody fills a request struct from a JSON or YAML file, or from stdin
// when path is "-".
//
// It exists because a request body with nested objects has no honest flag
// representation. Flags are applied after this, so a file can carry the bulk
// of a request and a flag can override one field of it.
func loadBody(path string, v any) error {
	data, err := readFileOrStdin(path)
	if err != nil {
		return err
	}
	trimmed := strings.TrimSpace(string(data))
	if trimmed == "" {
		return fmt.Errorf("--from-file %s is empty", path)
	}
	// JSON is a subset of YAML, so one decoder handles both. Trying JSON
	// first keeps its error messages, which point at a line and column.
	if strings.HasPrefix(trimmed, "{") || strings.HasPrefix(trimmed, "[") {
		if err := json.Unmarshal(data, v); err != nil {
			return fmt.Errorf("--from-file %s: %w", path, err)
		}
		return nil
	}
	// Round-trip through JSON so the struct's json tags apply: the SDK's
	// types carry json tags and no yaml ones.
	var generic any
	if err := yaml.Unmarshal(data, &generic); err != nil {
		return fmt.Errorf("--from-file %s: %w", path, err)
	}
	asJSON, err := json.Marshal(generic)
	if err != nil {
		return fmt.Errorf("--from-file %s: %w", path, err)
	}
	if err := json.Unmarshal(asJSON, v); err != nil {
		return fmt.Errorf("--from-file %s: %w", path, err)
	}
	return nil
}

// openBody returns a reader for an operation whose body is raw bytes,
// streaming from disk rather than buffering: an object upload can be larger
// than memory.
//
// The reader is wrapped in a progress reporter. This is the single chokepoint
// every streaming request body passes through, so one wrap covers every such
// operation the generator emits, now and later.
func openBody(path string) (io.Reader, func(), error) {
	if path == "" || path == "-" {
		// stdin has no length to be a fraction of, so the reporter shows
		// bytes and a rate rather than a percentage.
		r, done := progress.Wrap(os.Stdin, 0, "uploading")
		return r, done, nil
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	// A total makes the difference between "something is happening" and "this
	// finishes in nine minutes", which is what decides whether someone waits.
	var total int64
	if info, err := f.Stat(); err == nil {
		total = info.Size()
	}
	r, done := progress.Wrap(f, total, "uploading")
	return r, func() { done(); f.Close() }, nil
}

// readBody reads a whole body into memory, for the operations whose body is
// text the platform parses as one document.
func readBody(path string) (string, error) {
	data, err := readFileOrStdin(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func readFileOrStdin(path string) ([]byte, error) {
	if path == "" || path == "-" {
		return io.ReadAll(os.Stdin)
	}
	return os.ReadFile(path)
}

// parseTime accepts the timestamp forms a person is likely to type: RFC 3339,
// a bare date, or a Unix timestamp in seconds.
//
// The API takes RFC 3339, but requiring the full form on a command line means
// typing an hour and a timezone to ask for "since Monday".
func parseTime(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	for _, layout := range []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02",
	} {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	if secs, err := strconv.ParseInt(s, 10, 64); err == nil {
		return time.Unix(secs, 0).UTC(), nil
	}
	return time.Time{}, fmt.Errorf("cannot read %q as a time: use RFC 3339 (2026-01-15T09:30:00Z), a date (2026-01-15), or a Unix timestamp", s)
}
