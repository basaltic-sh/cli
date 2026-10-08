// Package progress reports how a long transfer is going.
//
// It exists because silence and a hang are indistinguishable. A 2.2GB upload
// took 28 minutes and printed nothing until it finished, and the natural
// response to a silent 28 minutes is Ctrl-C — which is the one action that
// used to leave a truncated object behind. A progress line is not only a
// comfort; it is what stops the dangerous response from being the reasonable
// one.
package progress

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/term"
)

// minimumSize is the smallest transfer worth drawing a bar for. Below this a
// progress line appears and vanishes in the same instant, which is noise.
const minimumSize = 2 << 20 // 2 MiB

// redrawInterval bounds how often the line is rewritten. Fast enough to look
// live, slow enough not to spend the transfer's time on formatting.
const redrawInterval = 200 * time.Millisecond

// Reader wraps an io.Reader and draws a progress line as it is consumed.
type Reader struct {
	inner io.Reader
	// total is the expected byte count, or 0 when it is not known — reading
	// from stdin, say, where there is nothing to be a fraction of.
	total   int64
	out     io.Writer
	label   string
	started time.Time

	mu       sync.Mutex
	read     int64
	lastDraw time.Time
	drawn    bool
	done     bool
}

// Wrap returns r wrapped in a progress reporter, or r unchanged when there is
// nobody to report to.
//
// Reporting goes to stderr and only when stderr is a terminal: a progress line
// in a pipeline is noise at best, and at worst lands in the middle of
// something being parsed. Same test the update notice uses.
func Wrap(r io.Reader, total int64, label string) (io.Reader, func()) {
	if !term.IsTerminal(int(os.Stderr.Fd())) {
		return r, func() {}
	}
	if total > 0 && total < minimumSize {
		return r, func() {}
	}
	return newReader(r, total, label, os.Stderr)
}

// newReader is Wrap without the terminal check, so the drawing can be tested
// against a buffer.
func newReader(r io.Reader, total int64, label string, out io.Writer) (*Reader, func()) {
	p := &Reader{inner: r, total: total, out: out, label: label, started: time.Now()}
	return p, p.Finish
}

// Size reports the total the reporter was given, or -1 when it has none.
//
// Load-bearing rather than convenience: the SDK sets Content-Length from a
// seekable body, and wrapping the file hides the *os.File from it. With no
// length declared the request goes out chunked, and a chunked upload cannot be
// checked for truncation server-side — so without this, adding a progress bar
// would silently disable the protection against a partial object.
//
// -1 AND NOT 0 for unknown. The SDK treats any non-negative answer as a real
// length, so returning 0 would declare an empty body and send nothing at all —
// turning "I do not know how big this is" into "this is empty". Reading from
// stdin is exactly that case.
func (p *Reader) Size() int64 {
	if p.total <= 0 {
		return -1
	}
	return p.total
}

func (p *Reader) Read(b []byte) (int, error) {
	n, err := p.inner.Read(b)
	if n > 0 {
		p.mu.Lock()
		p.read += int64(n)
		if time.Since(p.lastDraw) >= redrawInterval {
			p.draw()
			p.lastDraw = time.Now()
		}
		p.mu.Unlock()
	}
	return n, err
}

// Finish clears the progress line. Safe to call more than once.
//
// It clears rather than leaving a completed bar because the command prints its
// own result immediately afterwards, and two summaries of the same operation
// read as two operations.
func (p *Reader) Finish() {
	if p == nil {
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.done {
		return
	}
	p.done = true
	if p.drawn {
		// \r then spaces then \r: the line is overwritten rather than
		// scrolled, so nothing is left behind on a terminal that does not
		// support erase sequences.
		fmt.Fprintf(p.out, "\r%s\r", strings.Repeat(" ", 78))
	}
}

func (p *Reader) draw() {
	elapsed := time.Since(p.started)
	rate := float64(p.read) / elapsed.Seconds()

	var line string
	if p.total > 0 {
		pct := float64(p.read) / float64(p.total) * 100
		line = fmt.Sprintf("%s %s / %s  %5.1f%%  %s/s%s",
			p.label, Bytes(p.read), Bytes(p.total), pct, Bytes(int64(rate)), eta(p.read, p.total, rate))
	} else {
		// No total: still worth showing that bytes are moving, which is the
		// question a silent transfer raises.
		line = fmt.Sprintf("%s %s  %s/s", p.label, Bytes(p.read), Bytes(int64(rate)))
	}
	if len(line) > 78 {
		line = line[:78]
	}
	fmt.Fprintf(p.out, "\r%-78s", line)
	p.drawn = true
}

func eta(read, total int64, rate float64) string {
	if rate <= 0 || read >= total {
		return ""
	}
	remaining := time.Duration(float64(total-read)/rate) * time.Second
	if remaining > 99*time.Hour {
		return ""
	}
	return "  ETA " + remaining.Round(time.Second).String()
}

// Bytes renders a byte count the way a person reads one.
func Bytes(n int64) string {
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%d B", n)
	}
	div, exp := int64(unit), 0
	for v := n / unit; v >= unit; v /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(n)/float64(div), "KMGTPE"[exp])
}
