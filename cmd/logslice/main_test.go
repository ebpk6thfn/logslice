package main

import (
	"testing"
	"time"
)

func TestParseTime_RFC3339(t *testing.T) {
	input := "2024-03-15T10:30:00Z"
	got, err := parseTime(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := time.Date(2024, 3, 15, 10, 30, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestParseTime_ShortFormat(t *testing.T) {
	input := "2024-03-15T10:30:00"
	got, err := parseTime(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Year() != 2024 || got.Month() != 3 || got.Day() != 15 {
		t.Errorf("unexpected date components in %v", got)
	}
	if got.Hour() != 10 || got.Minute() != 30 || got.Second() != 0 {
		t.Errorf("unexpected time components in %v", got)
	}
}

func TestParseTime_InvalidString(t *testing.T) {
	input := "not-a-time"
	_, err := parseTime(input)
	if err == nil {
		t.Fatal("expected error for invalid time string, got nil")
	}
}

func TestParseTime_EmptyString(t *testing.T) {
	_, err := parseTime("")
	if err == nil {
		t.Fatal("expected error for empty string, got nil")
	}
}

func TestParseTime_RFC3339WithOffset(t *testing.T) {
	input := "2024-06-01T08:00:00+02:00"
	got, err := parseTime(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Convert to UTC for comparison
	gotUTC := got.UTC()
	want := time.Date(2024, 6, 1, 6, 0, 0, 0, time.UTC)
	if !gotUTC.Equal(want) {
		t.Errorf("got %v, want %v", gotUTC, want)
	}
}
