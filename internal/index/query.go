package index

import (
	"time"
)

// Entry represents a single indexed position in a log file.
type Entry struct {
	Offset    int64
	Timestamp time.Time
	Line      int
}

// Index holds a sorted slice of entries built from a log file.
type Index struct {
	entries []Entry
}

// Len returns the number of entries in the index.
func (idx *Index) Len() int {
	return len(idx.entries)
}

// FindStart returns the byte offset of the first entry whose timestamp
// is >= start. If all entries are before start, it returns -1.
func (idx *Index) FindStart(start time.Time) int64 {
	for _, e := range idx.entries {
		if !e.Timestamp.Before(start) {
			return e.Offset
		}
	}
	return -1
}

// FindEnd returns the byte offset of the first entry whose timestamp
// is > end. If no such entry exists, it returns -1 (meaning read to EOF).
func (idx *Index) FindEnd(end time.Time) int64 {
	for _, e := range idx.entries {
		if e.Timestamp.After(end) {
			return e.Offset
		}
	}
	return -1
}

// Entries returns a copy of all index entries.
func (idx *Index) Entries() []Entry {
	out := make([]Entry, len(idx.entries))
	copy(out, idx.entries)
	return out
}

// newIndex constructs an Index from a slice of entries.
func newIndex(entries []Entry) *Index {
	return &Index{entries: entries}
}
