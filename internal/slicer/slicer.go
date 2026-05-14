// Package slicer extracts time-bounded segments from log files
// using binary search via the seeker and line-by-line scanning.
package slicer

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/yourorg/logslice/internal/scanner"
)

// Options configures a Slice operation.
type Options struct {
	// TimestampFormat is an optional strftime-style or Go time layout.
	// If empty, auto-detection is attempted.
	TimestampFormat string

	// Start is the inclusive lower bound. Zero value means no lower bound.
	Start time.Time

	// End is the exclusive upper bound. Zero value means no upper bound.
	End time.Time
}

// Slice reads lines from the log file at path that fall within the
// time window defined by opts, writing matching lines to w.
// It uses binary search to seek to the start position, avoiding a
// full file scan when possible.
func Slice(path string, opts Options, w io.Writer) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("slicer: open %q: %w", path, err)
	}
	defer f.Close()

	seek, err := scanner.NewSeeker(f, opts.TimestampFormat)
	if err != nil {
		return fmt.Errorf("slicer: build seeker: %w", err)
	}

	var startOffset int64
	if !opts.Start.IsZero() {
		startOffset, err = seek.FindStart(opts.Start)
		if err != nil {
			return fmt.Errorf("slicer: find start: %w", err)
		}
	}

	if _, err = f.Seek(startOffset, io.SeekStart); err != nil {
		return fmt.Errorf("slicer: seek to offset %d: %w", startOffset, err)
	}

	sc := scanner.New(f, opts.TimestampFormat)
	for sc.Scan() {
		line := sc.Line()
		ts := sc.Timestamp()

		// Skip lines before the start window (can occur near binary-search boundary).
		if !opts.Start.IsZero() && !ts.IsZero() && ts.Before(opts.Start) {
			continue
		}

		// Stop once we exceed the end window.
		if !opts.End.IsZero() && !ts.IsZero() && !ts.Before(opts.End) {
			break
		}

		if _, err = fmt.Fprintln(w, line); err != nil {
			return fmt.Errorf("slicer: write: %w", err)
		}
	}

	return sc.Err()
}
