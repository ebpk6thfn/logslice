package index_test

import (
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/index"
)

const sampleLog = `2024-01-01T10:00:00Z INFO  starting up
2024-01-01T10:01:00Z DEBUG request received
2024-01-01T10:02:00Z INFO  processed ok
2024-01-01T10:03:00Z ERROR something failed
`

func mustTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

func buildIndex(t *testing.T, log string) index.Index {
	t.Helper()
	b := index.NewBuilder("")
	idx, err := b.Build(strings.NewReader(log))
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	return idx
}

func TestBuild_EntryCount(t *testing.T) {
	idx := buildIndex(t, sampleLog)
	if len(idx) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(idx))
	}
}

func TestBuild_FirstOffset(t *testing.T) {
	idx := buildIndex(t, sampleLog)
	if idx[0].Offset != 0 {
		t.Errorf("first offset should be 0, got %d", idx[0].Offset)
	}
}

func TestBuild_TimestampParsed(t *testing.T) {
	idx := buildIndex(t, sampleLog)
	want := mustTime("2024-01-01T10:01:00Z")
	if !idx[1].Timestamp.Equal(want) {
		t.Errorf("entry[1] timestamp = %v, want %v", idx[1].Timestamp, want)
	}
}

func TestFindStart_MatchesExact(t *testing.T) {
	idx := buildIndex(t, sampleLog)
	offset := idx.FindStart(mustTime("2024-01-01T10:02:00Z"))
	if offset != idx[2].Offset {
		t.Errorf("FindStart offset = %d, want %d", offset, idx[2].Offset)
	}
}

func TestFindStart_BeforeAll(t *testing.T) {
	idx := buildIndex(t, sampleLog)
	offset := idx.FindStart(mustTime("2023-01-01T00:00:00Z"))
	if offset != 0 {
		t.Errorf("expected offset 0 for before-all target, got %d", offset)
	}
}

func TestFindEnd_AfterAll(t *testing.T) {
	idx := buildIndex(t, sampleLog)
	offset := idx.FindEnd(mustTime("2025-01-01T00:00:00Z"))
	if offset != -1 {
		t.Errorf("expected -1 for after-all target, got %d", offset)
	}
}

func TestBuild_EmptyInput(t *testing.T) {
	idx := buildIndex(t, "")
	if len(idx) != 0 {
		t.Errorf("expected empty index, got %d entries", len(idx))
	}
}
