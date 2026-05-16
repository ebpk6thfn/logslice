package merge_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/merge"
)

func makeSource(entries []merge.Entry) <-chan merge.Entry {
	ch := make(chan merge.Entry, len(entries))
	for _, e := range entries {
		ch <- e
	}
	close(ch)
	return ch
}

func ts(offset int) time.Time {
	return time.Date(2024, 1, 1, 0, 0, offset, 0, time.UTC)
}

func TestMerge_SingleSource(t *testing.T) {
	src := makeSource([]merge.Entry{
		{Timestamp: ts(1), Line: []byte("a")},
		{Timestamp: ts(3), Line: []byte("b")},
	})
	out := merge.New(src).Merge()
	var got []merge.Entry
	for e := range out {
		got = append(got, e)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
	if string(got[0].Line) != "a" || string(got[1].Line) != "b" {
		t.Errorf("unexpected order: %v", got)
	}
}

func TestMerge_TwoSourcesInterleaved(t *testing.T) {
	src1 := makeSource([]merge.Entry{
		{Timestamp: ts(1), Line: []byte("s1-t1")},
		{Timestamp: ts(3), Line: []byte("s1-t3")},
	})
	src2 := makeSource([]merge.Entry{
		{Timestamp: ts(2), Line: []byte("s2-t2")},
		{Timestamp: ts(4), Line: []byte("s2-t4")},
	})
	out := merge.New(src1, src2).Merge()
	var lines []string
	for e := range out {
		lines = append(lines, string(e.Line))
	}
	want := []string{"s1-t1", "s2-t2", "s1-t3", "s2-t4"}
	for i, w := range want {
		if lines[i] != w {
			t.Errorf("pos %d: want %q, got %q", i, w, lines[i])
		}
	}
}

func TestMerge_EmptySources(t *testing.T) {
	src := makeSource(nil)
	out := merge.New(src).Merge()
	var count int
	for range out {
		count++
	}
	if count != 0 {
		t.Errorf("expected 0 entries, got %d", count)
	}
}

func TestMerge_SourceIndexPreserved(t *testing.T) {
	src1 := makeSource([]merge.Entry{{Timestamp: ts(1), Line: []byte("a")}})
	src2 := makeSource([]merge.Entry{{Timestamp: ts(2), Line: []byte("b")}})
	out := merge.New(src1, src2).Merge()
	first := <-out
	if first.Source != 0 {
		t.Errorf("expected source 0, got %d", first.Source)
	}
	second := <-out
	if second.Source != 1 {
		t.Errorf("expected source 1, got %d", second.Source)
	}
}
