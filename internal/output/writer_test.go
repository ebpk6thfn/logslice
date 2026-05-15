package output_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/output"
)

func TestWriter_RawFormat(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(output.Options{Dest: &buf, Format: output.FormatRaw})

	lines := []string{"line one", "line two", "line three"}
	for _, l := range lines {
		if err := w.WriteLine(l); err != nil {
			t.Fatalf("WriteLine error: %v", err)
		}
	}
	if err := w.Flush(); err != nil {
		t.Fatalf("Flush error: %v", err)
	}

	got := buf.String()
	for _, l := range lines {
		if !strings.Contains(got, l) {
			t.Errorf("expected output to contain %q", l)
		}
	}
}

func TestWriter_NumberedFormat(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(output.Options{Dest: &buf, Format: output.FormatNumbered})

	_ = w.WriteLine("alpha")
	_ = w.WriteLine("beta")
	_ = w.Flush()

	got := buf.String()
	if !strings.Contains(got, "1\talpha") {
		t.Errorf("expected '1\\talpha' in output, got: %q", got)
	}
	if !strings.Contains(got, "2\tbeta") {
		t.Errorf("expected '2\\tbeta' in output, got: %q", got)
	}
}

func TestWriter_LinesWritten(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(output.Options{Dest: &buf})

	for i := 0; i < 5; i++ {
		_ = w.WriteLine("entry")
	}
	if got := w.LinesWritten(); got != 5 {
		t.Errorf("LinesWritten() = %d, want 5", got)
	}
}

func TestWriter_DefaultsToRaw(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(output.Options{Dest: &buf})
	_ = w.WriteLine("hello")
	_ = w.Flush()

	if got := buf.String(); got != "hello\n" {
		t.Errorf("expected 'hello\\n', got %q", got)
	}
}

func TestWriter_FlushEmpty(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(output.Options{Dest: &buf})
	if err := w.Flush(); err != nil {
		t.Errorf("Flush on empty writer returned error: %v", err)
	}
}
