package slicer_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/slicer"
)

func TestOptions_Validate_MissingFilePath(t *testing.T) {
	opts := slicer.Options{}
	if err := opts.Validate(); err == nil {
		t.Error("expected error for missing FilePath")
	}
}

func TestOptions_Validate_EndBeforeStart(t *testing.T) {
	now := time.Now()
	opts := slicer.Options{
		FilePath: "some.log",
		Start:    now,
		End:      now.Add(-time.Minute),
	}
	err := opts.Validate()
	if err == nil {
		t.Error("expected error when End is before Start")
	}
}

func TestOptions_Validate_EndEqualsStart(t *testing.T) {
	now := time.Now()
	opts := slicer.Options{
		FilePath: "some.log",
		Start:    now,
		End:      now,
	}
	err := opts.Validate()
	if err == nil {
		t.Error("expected error when End equals Start")
	}
}

func TestOptions_Validate_Valid(t *testing.T) {
	now := time.Now()
	opts := slicer.Options{
		FilePath: "some.log",
		Start:    now,
		End:      now.Add(time.Hour),
	}
	if err := opts.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestOptions_Validate_ZeroTimes_Valid(t *testing.T) {
	opts := slicer.Options{
		FilePath: "some.log",
	}
	if err := opts.Validate(); err != nil {
		t.Errorf("unexpected error for zero times: %v", err)
	}
}

func TestOptionsError_Message(t *testing.T) {
	err := &slicer.OptionsError{Field: "FilePath", Msg: "must not be empty"}
	want := "logslice: options.FilePath: must not be empty"
	if err.Error() != want {
		t.Errorf("got %q, want %q", err.Error(), want)
	}
}
