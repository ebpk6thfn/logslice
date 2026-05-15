package stats_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/stats"
)

func TestCollector_InitialValues(t *testing.T) {
	c := stats.New()
	if c.LinesScanned != 0 || c.LinesMatched != 0 || c.LinesFiltered != 0 {
		t.Error("expected zero initial counters")
	}
	if c.StartTime.IsZero() {
		t.Error("expected StartTime to be set")
	}
}

func TestCollector_RecordLine(t *testing.T) {
	c := stats.New()
	c.RecordLine(true, false, 100)
	c.RecordLine(true, true, 50)
	c.RecordLine(false, false, 30)

	if c.LinesScanned != 3 {
		t.Errorf("LinesScanned: got %d, want 3", c.LinesScanned)
	}
	if c.LinesMatched != 2 {
		t.Errorf("LinesMatched: got %d, want 2", c.LinesMatched)
	}
	if c.LinesFiltered != 1 {
		t.Errorf("LinesFiltered: got %d, want 1", c.LinesFiltered)
	}
	if c.BytesRead != 180 {
		t.Errorf("BytesRead: got %d, want 180", c.BytesRead)
	}
}

func TestCollector_RecordWrite(t *testing.T) {
	c := stats.New()
	c.RecordWrite(512)
	c.RecordWrite(256)
	if c.BytesWritten != 768 {
		t.Errorf("BytesWritten: got %d, want 768", c.BytesWritten)
	}
}

func TestCollector_Elapsed(t *testing.T) {
	c := stats.New()
	time.Sleep(5 * time.Millisecond)
	c.Finish()
	if c.Elapsed() < 5*time.Millisecond {
		t.Errorf("elapsed too short: %s", c.Elapsed())
	}
	if c.EndTime.IsZero() {
		t.Error("EndTime should be set after Finish")
	}
}

func TestCollector_Elapsed_BeforeFinish(t *testing.T) {
	c := stats.New()
	time.Sleep(2 * time.Millisecond)
	if c.Elapsed() < 2*time.Millisecond {
		t.Errorf("elapsed before Finish too short: %s", c.Elapsed())
	}
}

func TestCollector_Print(t *testing.T) {
	c := stats.New()
	c.RecordLine(true, false, 200)
	c.RecordWrite(150)
	c.Finish()

	var buf bytes.Buffer
	c.Print(&buf)
	out := buf.String()

	for _, want := range []string{"Lines scanned", "Lines matched", "Bytes read", "Bytes written", "Elapsed"} {
		if !strings.Contains(out, want) {
			t.Errorf("Print output missing %q", want)
		}
	}
}
