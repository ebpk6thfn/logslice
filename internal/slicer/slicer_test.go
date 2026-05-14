package slicer_test

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/slicer"
)

func writeTempLog(t *testing.T, lines []string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "logslice-*.log")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	defer f.Close()
	for _, l := range lines {
		f.WriteString(l + "\n")
	}
	return f.Name()
}

var testLines = []string{
	"2024-01-10T10:00:00Z INFO starting service",
	"2024-01-10T10:01:00Z INFO request received",
	"2024-01-10T10:02:00Z WARN slow query detected",
	"2024-01-10T10:03:00Z ERROR connection timeout",
	"2024-01-10T10:04:00Z INFO request completed",
}

func mustParseRFC(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestSlice_WindowInMiddle(t *testing.T) {
	path := writeTempLog(t, testLines)
	var buf bytes.Buffer

	opts := slicer.Options{
		Start: mustParseRFC("2024-01-10T10:01:00Z"),
		End:   mustParseRFC("2024-01-10T10:03:00Z"),
	}
	if err := slicer.Slice(path, opts, &buf); err != nil {
		t.Fatalf("Slice error: %v", err)
	}

	got := strings.TrimSpace(buf.String())
	lines := strings.Split(got, "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d: %v", len(lines), lines)
	}
	if !strings.Contains(lines[0], "request received") {
		t.Errorf("unexpected first line: %s", lines[0])
	}
	if !strings.Contains(lines[1], "slow query") {
		t.Errorf("unexpected second line: %s", lines[1])
	}
}

func TestSlice_NoWindow_AllLines(t *testing.T) {
	path := writeTempLog(t, testLines)
	var buf bytes.Buffer

	if err := slicer.Slice(path, slicer.Options{}, &buf); err != nil {
		t.Fatalf("Slice error: %v", err)
	}

	got := strings.TrimSpace(buf.String())
	if count := len(strings.Split(got, "\n")); count != len(testLines) {
		t.Errorf("expected %d lines, got %d", len(testLines), count)
	}
}

func TestSlice_FileNotFound(t *testing.T) {
	var buf bytes.Buffer
	err := slicer.Slice("/nonexistent/path/file.log", slicer.Options{}, &buf)
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
