package scanner

import (
	"bufio"
	"io"
	"time"

	"github.com/user/logslice/internal/timestamp"
)

// Line represents a single log line with its parsed timestamp and raw content.
type Line struct {
	Timestamp time.Time
	Raw       string
	Offset    int64
}

// Scanner reads log lines from a reader and parses timestamps.
type Scanner struct {
	reader    *bufio.Reader
	format    string
	offset    int64
	current   Line
	err       error
}

// New creates a new Scanner with an optional explicit timestamp format.
// If format is empty, auto-detection is attempted.
func New(r io.Reader, format string) *Scanner {
	return &Scanner{
		reader: bufio.NewReaderSize(r, 64*1024),
		format: format,
	}
}

// Scan advances to the next log line. Returns true if a line was read.
func (s *Scanner) Scan() bool {
	lineStart := s.offset
	raw, err := s.reader.ReadString('\n')
	if len(raw) == 0 {
		s.err = err
		return false
	}
	s.offset += int64(len(raw))

	// Trim trailing newline for cleaner display but keep original for output.
	trimmed := raw
	if len(trimmed) > 0 && trimmed[len(trimmed)-1] == '\n' {
		trimmed = trimmed[:len(trimmed)-1]
	}

	ts, parseErr := timestamp.Parse(trimmed, s.format)
	if parseErr != nil {
		ts = time.Time{}
	}

	s.current = Line{
		Timestamp: ts,
		Raw:       raw,
		Offset:    lineStart,
	}
	return true
}

// Line returns the current line.
func (s *Scanner) Line() Line {
	return s.current
}

// Err returns the first non-EOF error encountered.
func (s *Scanner) Err() error {
	if s.err == io.EOF {
		return nil
	}
	return s.err
}
