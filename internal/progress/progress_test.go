package progress

import (
	"bytes"
	"io"
	"strings"
	"testing"
	"time"
)

func TestBytes(t *testing.T) {
	for in, want := range map[int64]string{
		0: "0 B", 512: "512 B", 1024: "1.0 KiB",
		1536: "1.5 KiB", 1048576: "1.0 MiB", 2323294720: "2.2 GiB",
	} {
		if got := Bytes(in); got != want {
			t.Errorf("Bytes(%d) = %q, want %q", in, got, want)
		}
	}
}

// The reporter must never change what the caller reads. It sits in the middle
// of an upload; corrupting a byte would be far worse than showing nothing.
func TestReaderPassesDataThroughUnchanged(t *testing.T) {
	payload := strings.Repeat("payload", 5000)
	var out bytes.Buffer
	r, done := newReader(strings.NewReader(payload), int64(len(payload)), "uploading", &out)
	got, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	done()
	if string(got) != payload {
		t.Errorf("the wrapped reader changed the data: %d bytes out, %d in", len(got), len(payload))
	}
}

func TestDrawsProgressWithATotal(t *testing.T) {
	var out bytes.Buffer
	p, done := newReader(strings.NewReader(strings.Repeat("x", 4096)), 4096, "uploading", &out)
	// Force a draw on the first read rather than waiting for the interval.
	p.lastDraw = time.Now().Add(-time.Second)
	buf := make([]byte, 1024)
	if _, err := p.Read(buf); err != nil {
		t.Fatal(err)
	}
	line := out.String()
	for _, want := range []string{"uploading", "KiB", "%", "/s"} {
		if !strings.Contains(line, want) {
			t.Errorf("progress line %q does not mention %q", strings.TrimSpace(line), want)
		}
	}
	if !strings.HasPrefix(line, "\r") {
		t.Error("the line does not start with a carriage return, so it will scroll rather than update in place")
	}
	done()
}

// Reading from stdin has no total, and a percentage of an unknown total is a
// lie. It must still show that bytes are moving, which is the question a
// silent transfer raises.
func TestDrawsBytesAndRateWithoutATotal(t *testing.T) {
	var out bytes.Buffer
	p, done := newReader(strings.NewReader(strings.Repeat("x", 4096)), 0, "uploading", &out)
	p.lastDraw = time.Now().Add(-time.Second)
	if _, err := p.Read(make([]byte, 1024)); err != nil {
		t.Fatal(err)
	}
	line := out.String()
	if strings.Contains(line, "%") {
		t.Errorf("a percentage was shown with no total to be a fraction of: %q", strings.TrimSpace(line))
	}
	if !strings.Contains(line, "/s") {
		t.Errorf("no rate shown: %q", strings.TrimSpace(line))
	}
	done()
}

// Finish clears the line rather than leaving a completed bar, because the
// command prints its own result immediately after and two summaries of one
// operation read as two operations.
func TestFinishClearsTheLineAndIsIdempotent(t *testing.T) {
	var out bytes.Buffer
	p, done := newReader(strings.NewReader("x"), 100, "uploading", &out)
	p.lastDraw = time.Now().Add(-time.Second)
	_, _ = p.Read(make([]byte, 1))
	out.Reset()

	done()
	cleared := out.String()
	if !strings.HasPrefix(cleared, "\r") || !strings.HasSuffix(cleared, "\r") {
		t.Errorf("Finish did not overwrite the line: %q", cleared)
	}
	out.Reset()
	done()
	if out.Len() != 0 {
		t.Errorf("a second Finish wrote %q; it must be a no-op", out.String())
	}
}

// Nothing was drawn, so nothing needs clearing — and writing a clear sequence
// would move the cursor over output the command had already printed.
func TestFinishWritesNothingIfNothingWasDrawn(t *testing.T) {
	var out bytes.Buffer
	_, done := newReader(strings.NewReader("x"), 100, "uploading", &out)
	done()
	if out.Len() != 0 {
		t.Errorf("Finish wrote %q without having drawn anything", out.String())
	}
}

// A transfer that finishes in an instant should not flash a progress line.
func TestWrapSkipsSmallTransfers(t *testing.T) {
	src := strings.NewReader("small")
	r, _ := Wrap(src, 1024, "uploading")
	if r != io.Reader(src) {
		t.Error("a small transfer was wrapped; the line would appear and vanish in the same instant")
	}
}

// In a pipeline there is nobody to show it to, and the line would land in the
// middle of something being parsed. Under `go test` stderr is not a terminal,
// so this is the path being exercised.
func TestWrapSkipsWhenStderrIsNotATerminal(t *testing.T) {
	src := strings.NewReader(strings.Repeat("x", 8<<20))
	r, _ := Wrap(src, 8<<20, "uploading")
	if r != io.Reader(src) {
		t.Error("the reader was wrapped despite stderr not being a terminal")
	}
}

// The SDK derives Content-Length from a seekable body, and wrapping hides the
// *os.File from it. Reporting Size keeps the length visible — without it,
// adding a progress bar would silently turn a length-declared upload into a
// chunked one, which cannot be checked for truncation server-side.
func TestReaderReportsItsSizeThroughTheWrapper(t *testing.T) {
	p, _ := newReader(strings.NewReader("payload"), 4096, "uploading", &bytes.Buffer{})
	var sized interface{ Size() int64 } = p
	if got := sized.Size(); got != 4096 {
		t.Errorf("Size() = %d, want 4096", got)
	}
}

// An unknown total must report -1, not 0. The SDK treats any non-negative
// answer as a real length, so 0 would declare an empty body and send nothing —
// silently uploading a zero-byte object instead of the piped data.
func TestReaderReportsUnknownSizeAsNegative(t *testing.T) {
	p, _ := newReader(strings.NewReader("payload"), 0, "uploading", &bytes.Buffer{})
	if got := p.Size(); got != -1 {
		t.Errorf("Size() = %d, want -1 for an unknown total; 0 would declare an empty body", got)
	}
}
