// Package progress provides a simple progress reporter for tracking
// how much of a log file has been processed during a slice operation.
package progress

import (
	"fmt"
	"io"
	"sync/atomic"
	"time"
)

// Reporter tracks and reports file processing progress.
type Reporter struct {
	total     int64
	processed atomic.Int64
	writer    io.Writer
	ticker    *time.Ticker
	done      chan struct{}
}

// New creates a new Reporter that writes progress updates to w.
// total is the total number of bytes in the file being processed.
// interval controls how often progress is reported; pass 0 to disable
// periodic reporting (updates are still available via Percent).
func New(w io.Writer, total int64, interval time.Duration) *Reporter {
	r := &Reporter{
		total:  total,
		writer: w,
		done:   make(chan struct{}),
	}
	if interval > 0 && total > 0 {
		r.ticker = time.NewTicker(interval)
		go r.run()
	}
	return r
}

// Advance records that n additional bytes have been processed.
func (r *Reporter) Advance(n int64) {
	r.processed.Add(n)
}

// Percent returns the current completion percentage (0–100).
// Returns 0 if total is unknown (zero).
func (r *Reporter) Percent() float64 {
	if r.total <= 0 {
		return 0
	}
	p := float64(r.processed.Load()) / float64(r.total) * 100
	if p > 100 {
		return 100
	}
	return p
}

// Stop halts periodic reporting and flushes a final 100% line when done is true.
func (r *Reporter) Stop(completed bool) {
	if r.ticker != nil {
		r.ticker.Stop()
		close(r.done)
	}
	if completed && r.writer != nil {
		fmt.Fprintf(r.writer, "progress: 100.00%%\n")
	}
}

func (r *Reporter) run() {
	for {
		select {
		case <-r.ticker.C:
			fmt.Fprintf(r.writer, "progress: %.2f%%\n", r.Percent())
		case <-r.done:
			return
		}
	}
}
