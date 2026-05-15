package index

import (
	"sort"
	"time"
)

// Entry represents a single indexed position within a log file.
type Entry struct {
	Offset    int64
	Timestamp time.Time
}

// index holds a sorted slice of entries and supports efficient time-based lookups.
type index struct {
	entries []Entry
}

// newIndex constructs an index from the given entries, sorting them by timestamp.
func newIndex(entries []Entry) *index {
	sorted := make([]Entry, len(entries))
	copy(sorted, entries)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Timestamp.Before(sorted[j].Timestamp)
	})
	return &index{entries: sorted}
}

// Len returns the number of entries in the index.
func (ix *index) Len() int {
	return len(ix.entries)
}

// FindStart returns the byte offset of the first entry whose timestamp is
// greater than or equal to t. If all entries precede t, it returns -1.
func (ix *index) FindStart(t time.Time) int64 {
	pos := sort.Search(len(ix.entries), func(i int) bool {
		return !ix.entries[i].Timestamp.Before(t)
	})
	if pos >= len(ix.entries) {
		return -1
	}
	return ix.entries[pos].Offset
}

// FindEnd returns the byte offset of the first entry whose timestamp is
// strictly greater than t. If no such entry exists, it returns -1.
func (ix *index) FindEnd(t time.Time) int64 {
	pos := sort.Search(len(ix.entries), func(i int) bool {
		return ix.entries[i].Timestamp.After(t)
	})
	if pos >= len(ix.entries) {
		return -1
	}
	return ix.entries[pos].Offset
}

// Entries returns a copy of all entries in sorted order.
func (ix *index) Entries() []Entry {
	out := make([]Entry, len(ix.entries))
	copy(out, ix.entries)
	return out
}
