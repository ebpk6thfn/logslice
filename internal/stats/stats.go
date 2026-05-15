// Package stats provides collection and reporting of slice operation metrics.
package stats

import (
	"fmt"
	"io"
	"time"
)

// Collector accumulates metrics during a log slicing operation.
type Collector struct {
	LinesScanned  int
	LinesMatched  int
	LinesFiltered int
	BytesRead     int64
	BytesWritten  int64
	StartTime     time.Time
	EndTime       time.Time
}

// New returns a new Collector with the start time set to now.
func New() *Collector {
	return &Collector{StartTime: time.Now()}
}

// Finish marks the end time of the operation.
func (c *Collector) Finish() {
	c.EndTime = time.Now()
}

// Elapsed returns the duration of the operation.
// If Finish has not been called, it returns the duration since StartTime.
func (c *Collector) Elapsed() time.Duration {
	if c.EndTime.IsZero() {
		return time.Since(c.StartTime)
	}
	return c.EndTime.Sub(c.StartTime)
}

// RecordLine records a scanned line, tracking whether it was matched or filtered.
func (c *Collector) RecordLine(matched, filtered bool, bytes int64) {
	c.LinesScanned++
	c.BytesRead += bytes
	if matched {
		c.LinesMatched++
	}
	if filtered {
		c.LinesFiltered++
	}
}

// RecordWrite records bytes written to output.
func (c *Collector) RecordWrite(bytes int64) {
	c.BytesWritten += bytes
}

// Print writes a human-readable summary of the collected stats to w.
func (c *Collector) Print(w io.Writer) {
	fmt.Fprintf(w, "Lines scanned : %d\n", c.LinesScanned)
	fmt.Fprintf(w, "Lines matched : %d\n", c.LinesMatched)
	fmt.Fprintf(w, "Lines filtered: %d\n", c.LinesFiltered)
	fmt.Fprintf(w, "Bytes read    : %d\n", c.BytesRead)
	fmt.Fprintf(w, "Bytes written : %d\n", c.BytesWritten)
	fmt.Fprintf(w, "Elapsed       : %s\n", c.Elapsed().Round(time.Millisecond))
}
