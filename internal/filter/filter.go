// Package filter provides line-level filtering for log entries.
package filter

import (
	"regexp"
	"strings"
)

// Options configures which log lines are included in output.
type Options struct {
	// Level filters lines by log level (e.g. "ERROR", "WARN"). Empty means all levels.
	Level string
	// Pattern is an optional regex that must match for a line to be included.
	Pattern string
}

// Filter decides whether log lines should be included in output.
type Filter struct {
	level   string
	pattern *regexp.Regexp
}

// New creates a Filter from the given Options.
// Returns an error if the regex pattern is invalid.
func New(opts Options) (*Filter, error) {
	f := &Filter{
		level: strings.ToUpper(opts.Level),
	}
	if opts.Pattern != "" {
		re, err := regexp.Compile(opts.Pattern)
		if err != nil {
			return nil, err
		}
		f.pattern = re
	}
	return f, nil
}

// Match reports whether the given log line passes all active filters.
func (f *Filter) Match(line string) bool {
	if f.level != "" {
		if !strings.Contains(strings.ToUpper(line), f.level) {
			return false
		}
	}
	if f.pattern != nil {
		if !f.pattern.MatchString(line) {
			return false
		}
	}
	return true
}

// IsNoop returns true when the filter will accept every line.
func (f *Filter) IsNoop() bool {
	return f.level == "" && f.pattern == nil
}
