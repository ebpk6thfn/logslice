package tail

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeLineToFile(t *testing.T, f *os.File, line string) {
	t.Helper()
	_, err := f.WriteString(line + "\n")
	if err != nil {
		t.Fatalf("write line: %v", err)
	}
}

func TestRotatedTailer_ReceivesLines(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "app.log")

	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create file: %v", err)
	}
	defer f.Close()

	rt, err := NewRotated(path, 20*time.Millisecond)
	if err != nil {
		t.Fatalf("NewRotated: %v", err)
	}
	defer rt.Stop()

	writeLineToFile(t, f, "line one")
	f.Sync()

	select {
	case got := <-rt.Lines():
		if got != "line one" {
			t.Errorf("expected 'line one', got %q", got)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timeout waiting for line")
	}
}

func TestRotatedTailer_FileNotFound(t *testing.T) {
	_, err := NewRotated("/nonexistent/path/app.log", 20*time.Millisecond)
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestRotatedTailer_DetectsRotation(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "app.log")

	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create file: %v", err)
	}

	rt, err := NewRotated(path, 20*time.Millisecond)
	if err != nil {
		t.Fatalf("NewRotated: %v", err)
	}
	defer rt.Stop()

	writeLineToFile(t, f, "before rotation")
	f.Sync()
	f.Close()

	select {
	case got := <-rt.Lines():
		if got != "before rotation" {
			t.Errorf("expected 'before rotation', got %q", got)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timeout waiting for pre-rotation line")
	}

	// Simulate rotation: remove old file, create new one
	os.Remove(path)
	f2, err := os.Create(path)
	if err != nil {
		t.Fatalf("create rotated file: %v", err)
	}
	defer f2.Close()

	writeLineToFile(t, f2, "after rotation")
	f2.Sync()

	select {
	case got := <-rt.Lines():
		if got != "after rotation" {
			t.Errorf("expected 'after rotation', got %q", got)
		}
	case <-time.After(800 * time.Millisecond):
		t.Fatal("timeout waiting for post-rotation line")
	}
}
