package timestamp

import (
	"fmt"
	"time"
)

// CommonFormats lists timestamp formats tried in order when auto-detecting.
var CommonFormats = []string{
	time.RFC3339Nano,
	time.RFC3339,
	"2006-01-02T15:04:05.999999999",
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05.999999999",
	"2006-01-02 15:04:05",
	"02/Jan/2006:15:04:05 -0700",
	"Jan 02 15:04:05",
	"Jan  2 15:04:05",
}

// Parse attempts to parse a timestamp string using a known format.
// If format is empty, it tries CommonFormats in order.
// Returns the parsed time and the matched format, or an error.
func Parse(s, format string) (time.Time, string, error) {
	if s == "" {
		return time.Time{}, "", fmt.Errorf("empty timestamp string")
	}

	if format != "" {
		t, err := time.Parse(format, s)
		if err != nil {
			return time.Time{}, "", fmt.Errorf("parse with format %q: %w", format, err)
		}
		return t, format, nil
	}

	for _, f := range CommonFormats {
		t, err := time.Parse(f, s)
		if err == nil {
			return t, f, nil
		}
	}

	return time.Time{}, "", fmt.Errorf("could not parse timestamp %q with any known format", s)
}

// MustParse is like Parse but panics on error. Useful in tests.
func MustParse(s, format string) time.Time {
	t, _, err := Parse(s, format)
	if err != nil {
		panic(err)
	}
	return t
}
