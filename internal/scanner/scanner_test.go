package scanner

import (
	"strings"
	"testing"
	"time"
)

const sampleLogs = `2024-01-15T10:00:00Z INFO starting service
2024-01-15T10:01:00Z DEBUG connected to database
2024-01-15T10:02:00Z WARN high memory usage
2024-01-15T10:03:00Z ERROR connection timeout
`

func TestScanner_BasicScan(t *testing.T) {
	r := strings.NewReader(sampleLogs)
	s := New(r, "")

	var lines []Line
	for s.Scan() {
		lines = append(lines, s.Line())
	}
	if err := s.Err(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 4 {
		t.Fatalf("expected 4 lines, got %d", len(lines))
	}
}

func TestScanner_TimestampParsed(t *testing.T) {
	r := strings.NewReader(sampleLogs)
	s := New(r, time.RFC3339)

	if !s.Scan() {
		t.Fatal("expected at least one line")
	}
	line := s.Line()
	if line.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp on first line")
	}
	expected := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	if !line.Timestamp.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, line.Timestamp)
	}
}

func TestScanner_OffsetTracking(t *testing.T) {
	r := strings.NewReader(sampleLogs)
	s := New(r, "")

	var offsets []int64
	for s.Scan() {
		offsets = append(offsets, s.Line().Offset)
	}
	if len(offsets) == 0 {
		t.Fatal("no lines scanned")
	}
	if offsets[0] != 0 {
		t.Errorf("first offset should be 0, got %d", offsets[0])
	}
	for i := 1; i < len(offsets); i++ {
		if offsets[i] <= offsets[i-1] {
			t.Errorf("offsets should be strictly increasing: %v", offsets)
		}
	}
}

func TestScanner_EmptyInput(t *testing.T) {
	r := strings.NewReader("")
	s := New(r, "")
	if s.Scan() {
		t.Error("expected no lines from empty input")
	}
	if s.Err() != nil {
		t.Errorf("unexpected error: %v", s.Err())
	}
}

func TestScanner_NoTimestampLine(t *testing.T) {
	r := strings.NewReader("this line has no timestamp\n")
	s := New(r, "")
	if !s.Scan() {
		t.Fatal("expected one line")
	}
	if !s.Line().Timestamp.IsZero() {
		t.Error("expected zero timestamp for unparseable line")
	}
}
