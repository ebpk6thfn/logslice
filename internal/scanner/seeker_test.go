package scanner

import (
	"strings"
	"testing"
	"time"
)

func makeSeeker(content string) *Seeker {
	return NewSeeker(strings.NewReader(content), time.RFC3339)
}

const seekerLogs = "2024-03-01T08:00:00Z level=info msg=boot\n" +
	"2024-03-01T08:05:00Z level=info msg=ready\n" +
	"2024-03-01T08:10:00Z level=warn msg=slow\n" +
	"2024-03-01T08:15:00Z level=error msg=down\n"

func TestSeeker_FindStart(t *testing.T) {
	sk := makeSeeker(seekerLogs)
	target := time.Date(2024, 3, 1, 8, 5, 0, 0, time.UTC)
	offset, err := sk.FindStart(0, target)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if offset <= 0 {
		t.Errorf("expected positive offset, got %d", offset)
	}
	// The first line starts at 0; second line should be at a positive offset.
	firstLineLen := int64(len("2024-03-01T08:00:00Z level=info msg=boot\n"))
	if offset != firstLineLen {
		t.Errorf("expected offset %d, got %d", firstLineLen, offset)
	}
}

func TestSeeker_FindStart_BeforeAll(t *testing.T) {
	sk := makeSeeker(seekerLogs)
	target := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	offset, err := sk.FindStart(0, target)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if offset != 0 {
		t.Errorf("expected offset 0 for target before all entries, got %d", offset)
	}
}

func TestSeeker_FindStart_AfterAll(t *testing.T) {
	sk := makeSeeker(seekerLogs)
	target := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	offset, err := sk.FindStart(0, target)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if offset != -1 {
		t.Errorf("expected -1 when target is after all entries, got %d", offset)
	}
}

func TestSeeker_FindEnd(t *testing.T) {
	sk := makeSeeker(seekerLogs)
	target := time.Date(2024, 3, 1, 8, 10, 0, 0, time.UTC)
	endOffset, err := sk.FindEnd(0, target)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := int64(len("2024-03-01T08:00:00Z level=info msg=boot\n") +
		len("2024-03-01T08:05:00Z level=info msg=ready\n") +
		len("2024-03-01T08:10:00Z level=warn msg=slow\n"))
	if endOffset != expected {
		t.Errorf("expected end offset %d, got %d", expected, endOffset)
	}
}
