package index

import (
	"testing"
	"time"
)

func makeEntries(timestamps ...string) []Entry {
	entries := make([]Entry, len(timestamps))
	for i, ts := range timestamps {
		t, err := time.Parse(time.RFC3339, ts)
		if err != nil {
			panic(err)
		}
		entries[i] = Entry{Offset: int64(i * 100), Timestamp: t}
	}
	return entries
}

func TestIndex_Len(t *testing.T) {
	ix := newIndex(makeEntries(
		"2024-01-01T00:00:00Z",
		"2024-01-01T01:00:00Z",
	))
	if got := ix.Len(); got != 2 {
		t.Fatalf("Len() = %d, want 2", got)
	}
}

func TestIndex_FindStart_Exact(t *testing.T) {
	ix := newIndex(makeEntries(
		"2024-01-01T00:00:00Z",
		"2024-01-01T01:00:00Z",
		"2024-01-01T02:00:00Z",
	))
	target, _ := time.Parse(time.RFC3339, "2024-01-01T01:00:00Z")
	got := ix.FindStart(target)
	if got != 100 {
		t.Fatalf("FindStart exact = %d, want 100", got)
	}
}

func TestIndex_FindStart_BeforeAll(t *testing.T) {
	ix := newIndex(makeEntries(
		"2024-01-01T01:00:00Z",
		"2024-01-01T02:00:00Z",
	))
	target, _ := time.Parse(time.RFC3339, "2024-01-01T00:00:00Z")
	got := ix.FindStart(target)
	if got != 0 {
		t.Fatalf("FindStart before all = %d, want 0", got)
	}
}

func TestIndex_FindStart_AfterAll(t *testing.T) {
	ix := newIndex(makeEntries(
		"2024-01-01T00:00:00Z",
		"2024-01-01T01:00:00Z",
	))
	target, _ := time.Parse(time.RFC3339, "2024-01-01T03:00:00Z")
	got := ix.FindStart(target)
	if got != -1 {
		t.Fatalf("FindStart after all = %d, want -1", got)
	}
}

func TestIndex_FindEnd_Exclusive(t *testing.T) {
	ix := newIndex(makeEntries(
		"2024-01-01T00:00:00Z",
		"2024-01-01T01:00:00Z",
		"2024-01-01T02:00:00Z",
	))
	target, _ := time.Parse(time.RFC3339, "2024-01-01T01:00:00Z")
	got := ix.FindEnd(target)
	if got != 200 {
		t.Fatalf("FindEnd exclusive = %d, want 200", got)
	}
}

func TestIndex_FindEnd_AfterAll(t *testing.T) {
	ix := newIndex(makeEntries(
		"2024-01-01T00:00:00Z",
		"2024-01-01T01:00:00Z",
	))
	target, _ := time.Parse(time.RFC3339, "2024-01-01T05:00:00Z")
	got := ix.FindEnd(target)
	if got != -1 {
		t.Fatalf("FindEnd after all = %d, want -1", got)
	}
}

func TestIndex_Entries_ReturnsCopy(t *testing.T) {
	ix := newIndex(makeEntries(
		"2024-01-01T00:00:00Z",
	))
	e := ix.Entries()
	e[0].Offset = 9999
	if ix.entries[0].Offset == 9999 {
		t.Fatal("Entries() returned a reference to internal slice")
	}
}
