package timestamp_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/timestamp"
)

func TestParse_ExplicitFormat(t *testing.T) {
	s := "2024-03-15T10:22:33Z"
	got, fmt, err := timestamp.Parse(s, time.RFC3339)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fmt != time.RFC3339 {
		t.Errorf("expected format %q, got %q", time.RFC3339, fmt)
	}
	if got.Year() != 2024 || got.Month() != 3 || got.Day() != 15 {
		t.Errorf("unexpected time value: %v", got)
	}
}

func TestParse_AutoDetect(t *testing.T) {
	cases := []struct {
		input      string
		wantFormat string
	}{
		{"2024-03-15T10:22:33Z", time.RFC3339},
		{"2024-03-15 10:22:33", "2006-01-02 15:04:05"},
		{"2024-03-15T10:22:33.123456789", "2006-01-02T15:04:05.999999999"},
		{"15/Mar/2024:10:22:33 +0000", "02/Jan/2006:15:04:05 -0700"},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			_, gotFmt, err := timestamp.Parse(tc.input, "")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if gotFmt != tc.wantFormat {
				t.Errorf("expected format %q, got %q", tc.wantFormat, gotFmt)
			}
		})
	}
}

func TestParse_EmptyString(t *testing.T) {
	_, _, err := timestamp.Parse("", "")
	if err == nil {
		t.Fatal("expected error for empty string, got nil")
	}
}

func TestParse_UnknownFormat(t *testing.T) {
	_, _, err := timestamp.Parse("not-a-timestamp", "")
	if err == nil {
		t.Fatal("expected error for unknown format, got nil")
	}
}

func TestMustParse_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic, got none")
		}
	}()
	timestamp.MustParse("bad input", "")
}
