package progress_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/progress"
)

func TestReporter_Percent_Zero(t *testing.T) {
	r := progress.New(nil, 0, 0)
	if got := r.Percent(); got != 0 {
		t.Errorf("expected 0 for unknown total, got %f", got)
	}
}

func TestReporter_Percent_Advances(t *testing.T) {
	r := progress.New(nil, 200, 0)
	r.Advance(50)
	if got := r.Percent(); got != 25.0 {
		t.Errorf("expected 25.0, got %f", got)
	}
	r.Advance(150)
	if got := r.Percent(); got != 100.0 {
		t.Errorf("expected 100.0, got %f", got)
	}
}

func TestReporter_Percent_Clamps(t *testing.T) {
	r := progress.New(nil, 100, 0)
	r.Advance(200)
	if got := r.Percent(); got != 100.0 {
		t.Errorf("expected clamped 100.0, got %f", got)
	}
}

func TestReporter_Stop_WritesCompletion(t *testing.T) {
	var buf bytes.Buffer
	r := progress.New(&buf, 100, 0)
	r.Advance(100)
	r.Stop(true)
	if !strings.Contains(buf.String(), "100.00%") {
		t.Errorf("expected completion line in output, got: %q", buf.String())
	}
}

func TestReporter_Stop_NoWriteWhenNotCompleted(t *testing.T) {
	var buf bytes.Buffer
	r := progress.New(&buf, 100, 0)
	r.Advance(50)
	r.Stop(false)
	if buf.Len() != 0 {
		t.Errorf("expected no output when not completed, got: %q", buf.String())
	}
}

func TestReporter_PeriodicReporting(t *testing.T) {
	var buf bytes.Buffer
	r := progress.New(&buf, 1000, 20*time.Millisecond)
	r.Advance(500)
	time.Sleep(60 * time.Millisecond)
	r.Stop(false)
	if !strings.Contains(buf.String(), "progress:") {
		t.Errorf("expected periodic progress lines, got: %q", buf.String())
	}
}
