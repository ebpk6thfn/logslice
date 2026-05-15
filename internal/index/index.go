// Package index provides a lightweight line-offset index for large log files,
// enabling fast binary-search-based seeking without full file scans.
package index

import (
	"bufio"
	"io"
	"time"

	"github.com/yourorg/logslice/internal/timestamp"
)

// Entry records the byte offset and parsed timestamp of a single log line.
type Entry struct {
	Offset    int64
	Timestamp time.Time
}

// Index is an ordered slice of Entry values built from a log file.
type Index []Entry

// Builder constructs an Index by scanning a ReadSeeker.
type Builder struct {
	format string
}

// NewBuilder returns a Builder that uses the given timestamp format.
// Pass an empty string to enable auto-detection.
func NewBuilder(format string) *Builder {
	return &Builder{format: format}
}

// Build reads from r, recording the offset and timestamp of every line that
// contains a parseable timestamp. Lines without a timestamp are skipped.
func (b *Builder) Build(r io.ReadSeeker) (Index, error) {
	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	var idx Index
	var offset int64

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		lineLen := int64(len(line)) + 1 // +1 for newline

		ts, err := timestamp.Parse(line, b.format)
		if err == nil {
			idx = append(idx, Entry{Offset: offset, Timestamp: ts})
		}

		offset += lineLen
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return idx, nil
}

// FindStart returns the byte offset of the first entry whose timestamp is
// >= target. Returns 0 if all entries are after target or the index is empty.
func (idx Index) FindStart(target time.Time) int64 {
	for _, e := range idx {
		if !e.Timestamp.Before(target) {
			return e.Offset
		}
	}
	return 0
}

// FindEnd returns the byte offset of the first entry whose timestamp is
// >= target. Returns -1 to signal "read to EOF" when no such entry exists.
func (idx Index) FindEnd(target time.Time) int64 {
	for _, e := range idx {
		if !e.Timestamp.Before(target) {
			return e.Offset
		}
	}
	return -1
}
