package scanner

import (
	"io"
	"time"
)

// Seeker finds byte offsets in a ReadSeeker where timestamps cross boundaries.
type Seeker struct {
	rs     io.ReadSeeker
	format string
}

// NewSeeker creates a Seeker for binary-search style offset finding.
func NewSeeker(rs io.ReadSeeker, format string) *Seeker {
	return &Seeker{rs: rs, format: format}
}

// FindStart returns the byte offset of the first line whose timestamp >= t.
// It performs a linear scan from the given startOffset.
// Returns -1 if no matching line is found.
func (sk *Seeker) FindStart(startOffset int64, t time.Time) (int64, error) {
	if _, err := sk.rs.Seek(startOffset, io.SeekStart); err != nil {
		return 0, err
	}
	s := New(sk.rs, sk.format)
	for s.Scan() {
		line := s.Line()
		if !line.Timestamp.IsZero() && !line.Timestamp.Before(t) {
			return line.Offset, nil
		}
	}
	if err := s.Err(); err != nil {
		return 0, err
	}
	return -1, nil
}

// FindEnd returns the byte offset just past the last line whose timestamp <= t.
// It performs a linear scan from the given startOffset.
// Returns startOffset if no lines with a timestamp <= t are found.
func (sk *Seeker) FindEnd(startOffset int64, t time.Time) (int64, error) {
	if _, err := sk.rs.Seek(startOffset, io.SeekStart); err != nil {
		return 0, err
	}
	s := New(sk.rs, sk.format)
	var lastEnd int64 = startOffset
	for s.Scan() {
		line := s.Line()
		if !line.Timestamp.IsZero() && line.Timestamp.After(t) {
			break
		}
		lastEnd = line.Offset + int64(len(line.Raw))
	}
	if err := s.Err(); err != nil {
		return 0, err
	}
	return lastEnd, nil
}

// FindRange returns the start and end byte offsets for lines whose timestamps
// fall within [from, to]. It is a convenience wrapper around FindStart and FindEnd.
func (sk *Seeker) FindRange(from, to time.Time) (start, end int64, err error) {
	start, err = sk.FindStart(0, from)
	if err != nil {
		return 0, 0, err
	}
	if start < 0 {
		return -1, -1, nil
	}
	end, err = sk.FindEnd(start, to)
	if err != nil {
		return 0, 0, err
	}
	return start, end, nil
}
