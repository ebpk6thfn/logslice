package index

import (
	"testing"
	"time"
)

func makeEntries() []Entry {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	return []Entry{
		{Offset: 0, Timestamp: base, Line: 1},
		{Offset: 100, Timestamp: base.Add(1 * time.Minute), Line: 5},
		{Offset: 250, Timestamp: base.Add(2 * time.Minute), Line: 10},
		{Offset: 400, Timestamp: base.Add(3 * time.Minute), Line: 15},
	}
}

func TestIndex_Len(t *testing.T) {
	idx := newIndex(makeEntries())
	if idx.Len() != 4 {
		t.Fatalf("expected 4, got %d", idx.Len())
	}
}

func TestIndex_FindStart_Exact(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	idx := newIndex(makeEntries())
	offset := idx.FindStart(base.Add(1 * time.Minute))
	if offset != 100 {
		t.Fatalf("expected 100, got %d", offset)
	}
}

func TestIndex_FindStart_BeforeAll(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	idx := newIndex(makeEntries())
	offset := idx.FindStart(base.Add(-5 * time.Minute))
	if offset != 0 {
		t.Fatalf("expected 0, got %d", offset)
	}
}

func TestIndex_FindStart_AfterAll(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	idx := newIndex(makeEntries())
	offset := idx.FindStart(base.Add(10 * time.Minute))
	if offset != -1 {
		t.Fatalf("expected -1, got %d", offset)
	}
}

func TestIndex_FindEnd_Middle(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	idx := newIndex(makeEntries())
	offset := idx.FindEnd(base.Add(1 * time.Minute))
	if offset != 250 {
		t.Fatalf("expected 250, got %d", offset)
	}
}

func TestIndex_FindEnd_AfterAll(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	idx := newIndex(makeEntries())
	offset := idx.FindEnd(base.Add(10 * time.Minute))
	if offset != -1 {
		t.Fatalf("expected -1 (read to EOF), got %d", offset)
	}
}

func TestIndex_Entries_ReturnsCopy(t *testing.T) {
	idx := newIndex(makeEntries())
	e := idx.Entries()
	e[0].Offset = 9999
	if idx.entries[0].Offset == 9999 {
		t.Fatal("Entries() should return a copy, not a reference")
	}
}
