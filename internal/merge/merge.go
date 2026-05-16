// Package merge provides utilities for merging multiple sorted log
// streams into a single chronologically ordered output.
package merge

import (
	"container/heap"
	"time"
)

// Entry represents a single log line with its parsed timestamp and source
// index, used during the merge process.
type Entry struct {
	Timestamp time.Time
	Line      []byte
	Source    int
}

// entryHeap implements heap.Interface for min-heap ordering by timestamp.
type entryHeap []Entry

func (h entryHeap) Len() int            { return len(h) }
func (h entryHeap) Less(i, j int) bool  { return h[i].Timestamp.Before(h[j].Timestamp) }
func (h entryHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *entryHeap) Push(x interface{}) { *h = append(*h, x.(Entry)) }
func (h *entryHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

// Merger merges pre-sorted Entry channels into a single ordered stream.
type Merger struct {
	sources []<-chan Entry
}

// New creates a Merger that will merge the provided source channels.
func New(sources ...(<-chan Entry)) *Merger {
	return &Merger{sources: sources}
}

// Merge reads from all source channels and emits entries in ascending
// timestamp order via the returned channel. The channel is closed once
// all sources are exhausted.
func (m *Merger) Merge() <-chan Entry {
	out := make(chan Entry, 64)
	go func() {
		defer close(out)
		h := &entryHeap{}
		heap.Init(h)

		// Seed the heap with the first entry from each source.
		for i, src := range m.sources {
			if e, ok := <-src; ok {
				e.Source = i
				heap.Push(h, e)
			}
		}

		for h.Len() > 0 {
			smallest := heap.Pop(h).(Entry)
			out <- smallest
			// Refill from the same source.
			if e, ok := <-m.sources[smallest.Source]; ok {
				e.Source = smallest.Source
				heap.Push(h, e)
			}
		}
	}()
	return out
}
