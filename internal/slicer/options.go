package slicer

import (
	"time"

	"github.com/yourorg/logslice/internal/filter"
)

// Options controls how Slice extracts log segments.
type Options struct {
	// FilePath is the path to the log file to slice.
	FilePath string
	// Start is the beginning of the time window (inclusive).
	Start time.Time
	// End is the end of the time window (exclusive).
	End time.Time
	// TimestampFormat is the strftime-style format used to parse timestamps.
	// Leave empty for auto-detection.
	TimestampFormat string
	// Filter optionally restricts which matching lines are emitted.
	Filter *filter.Filter
}

// Validate returns an error if the options are logically inconsistent.
func (o Options) Validate() error {
	if o.FilePath == "" {
		return &OptionsError{Field: "FilePath", Msg: "must not be empty"}
	}
	if !o.Start.IsZero() && !o.End.IsZero() && !o.End.After(o.Start) {
		return &OptionsError{Field: "End", Msg: "must be after Start"}
	}
	return nil
}

// OptionsError describes an invalid option.
type OptionsError struct {
	Field string
	Msg   string
}

func (e *OptionsError) Error() string {
	return "logslice: options." + e.Field + ": " + e.Msg
}
