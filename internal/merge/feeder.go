package merge

import (
	"bufio"
	"io"
	"time"

	"github.com/yourorg/logslice/internal/timestamp"
)

// FeedOptions controls how a feeder reads and parses lines.
type FeedOptions struct {
	// Format is the explicit timestamp format; empty means auto-detect.
	Format string
	// BufferSize is the channel buffer depth (default 64).
	BufferSize int
}

// Feed reads lines from r, parses a leading timestamp from each line using
// the given format (or auto-detection when format is empty), and sends
// Entry values to the returned channel. Lines that cannot be parsed receive
// the zero time and are still forwarded so callers can decide what to do.
func Feed(r io.Reader, opts FeedOptions) <-chan Entry {
	buf := opts.BufferSize
	if buf <= 0 {
		buf = 64
	}
	ch := make(chan Entry, buf)
	go func() {
		defer close(ch)
		sc := bufio.NewScanner(r)
		for sc.Scan() {
			line := sc.Bytes()
			copy := make([]byte, len(line))
			copy(copy, line)
			var t time.Time
			if opts.Format != "" {
				t, _ = timestamp.Parse(string(line), opts.Format)
			} else {
				t, _ = timestamp.Parse(string(line), "")
			}
			ch <- Entry{Timestamp: t, Line: copy}
		}
	}()
	return ch
}
