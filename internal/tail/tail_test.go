package tail_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/tail"
)

func writeLine(t *testing.T, f *os.File, s string) {
	t.Helper()
	_, err := f.WriteString(s + "\n")
	if err != nil {
		t.Fatalf("writeLine: %v", err)
	}
}

func TestTailer_ReceivesNewLines(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "tail-*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	// Write a line before tailing starts — should NOT appear.
	writeLine(t, f, "old line")

	tlr := tail.New(f.Name(), 50*time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- tlr.Follow(ctx)
	}()

	// Give the tailer time to seek to end.
	time.Sleep(100 * time.Millisecond)

	writeLine(t, f, "line one")
	writeLine(t, f, "line two")

	var received []string
	timeout := time.After(1 * time.Second)
collect:
	for {
		select {
		case l, ok := <-tlr.Lines():
			if !ok {
				break collect
			}
			received = append(received, l.Text)
			if len(received) == 2 {
				cancel()
			}
		case <-timeout:
			t.Fatal("timed out waiting for lines")
		}
	}

	if len(received) != 2 {
		t.Fatalf("expected 2 lines, got %d: %v", len(received), received)
	}
	if received[0] != "line one" {
		t.Errorf("line[0] = %q, want %q", received[0], "line one")
	}
	if received[1] != "line two" {
		t.Errorf("line[1] = %q, want %q", received[1], "line two")
	}
}

func TestTailer_FileNotFound(t *testing.T) {
	tlr := tail.New("/nonexistent/path/file.log", 50*time.Millisecond)
	err := tlr.Follow(context.Background())
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestTailer_DefaultPollInterval(t *testing.T) {
	// Passing 0 should not panic and should use the default interval.
	tlr := tail.New("/dev/null", 0)
	if tlr == nil {
		t.Fatal("New returned nil")
	}
}
