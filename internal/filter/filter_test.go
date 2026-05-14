package filter_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/filter"
)

func TestFilter_NoOptions_AcceptsAll(t *testing.T) {
	f, err := filter.New(filter.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.IsNoop() {
		t.Error("expected IsNoop to be true for empty options")
	}
	if !f.Match("2024-01-01 INFO some message") {
		t.Error("expected match for any line with no filters")
	}
}

func TestFilter_LevelFilter(t *testing.T) {
	f, err := filter.New(filter.Options{Level: "ERROR"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Match("2024-01-01 INFO hello world") {
		t.Error("INFO line should not match ERROR filter")
	}
	if !f.Match("2024-01-01 ERROR something failed") {
		t.Error("ERROR line should match ERROR filter")
	}
}

func TestFilter_LevelFilter_CaseInsensitive(t *testing.T) {
	f, err := filter.New(filter.Options{Level: "warn"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.Match("2024-01-01 WARN disk usage high") {
		t.Error("WARN line should match warn filter")
	}
}

func TestFilter_PatternFilter(t *testing.T) {
	f, err := filter.New(filter.Options{Pattern: `user=\d+`})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.Match("2024-01-01 INFO login user=42") {
		t.Error("line with user=42 should match pattern")
	}
	if f.Match("2024-01-01 INFO logout admin") {
		t.Error("line without user=<digits> should not match pattern")
	}
}

func TestFilter_InvalidPattern(t *testing.T) {
	_, err := filter.New(filter.Options{Pattern: `[invalid`})
	if err == nil {
		t.Error("expected error for invalid regex pattern")
	}
}

func TestFilter_LevelAndPattern_BothMustMatch(t *testing.T) {
	f, err := filter.New(filter.Options{Level: "ERROR", Pattern: "timeout"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Match("2024-01-01 ERROR connection refused") {
		t.Error("ERROR line without 'timeout' should not match")
	}
	if f.Match("2024-01-01 WARN timeout reached") {
		t.Error("WARN line with 'timeout' should not match")
	}
	if !f.Match("2024-01-01 ERROR timeout after 30s") {
		t.Error("ERROR line with 'timeout' should match")
	}
}
