package merge_test

import (
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/merge"
)

func TestFeed_ParsesTimestamps(t *testing.T) {
	input := strings.Join([]string{
		"2024-01-01T00:00:01Z line one",
		"2024-01-01T00:00:02Z line two",
	}, "\n")

	ch := merge.Feed(strings.NewReader(input), merge.FeedOptions{})
	var entries []merge.Entry
	for e := range ch {
		entries = append(entries, e)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Timestamp.IsZero() {
		t.Error("expected non-zero timestamp for first entry")
	}
	if !entries[0].Timestamp.Before(entries[1].Timestamp) {
		t.Error("expected first timestamp before second")
	}
}

func TestFeed_UnparsableLineGetsZeroTime(t *testing.T) {
	input := "this line has no timestamp at all"
	ch := merge.Feed(strings.NewReader(input), merge.FeedOptions{})
	e := <-ch
	if !e.Timestamp.IsZero() {
		t.Errorf("expected zero timestamp, got %v", e.Timestamp)
	}
	if string(e.Line) != input {
		t.Errorf("expected line %q, got %q", input, e.Line)
	}
}

func TestFeed_EmptyReader(t *testing.T) {
	ch := merge.Feed(strings.NewReader(""), merge.FeedOptions{})
	var count int
	for range ch {
		count++
	}
	if count != 0 {
		t.Errorf("expected 0 entries from empty reader, got %d", count)
	}
}

func TestFeed_CustomBufferSize(t *testing.T) {
	input := "2024-01-01T00:00:01Z only one line"
	ch := merge.Feed(strings.NewReader(input), merge.FeedOptions{BufferSize: 1})
	e := <-ch
	if string(e.Line) != input {
		t.Errorf("unexpected line: %q", e.Line)
	}
}
