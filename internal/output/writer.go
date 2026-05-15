// Package output handles writing sliced log segments to various destinations.
package output

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// Format represents the output format for log lines.
type Format string

const (
	// FormatRaw writes lines as-is, preserving original formatting.
	FormatRaw Format = "raw"
	// FormatNumbered prefixes each line with its line number.
	FormatNumbered Format = "numbered"
)

// Options configures the output writer.
type Options struct {
	// Dest is the writer to send output to. Defaults to os.Stdout if nil.
	Dest io.Writer
	// Format controls how lines are written.
	Format Format
	// LineBufferSize controls the size of the write buffer.
	LineBufferSize int
}

// Writer wraps an io.Writer with buffering and formatting support.
type Writer struct {
	bw      *bufio.Writer
	format  Format
	lineNum int
}

// New creates a new Writer with the given options.
func New(opts Options) *Writer {
	dest := opts.Dest
	if dest == nil {
		dest = os.Stdout
	}
	bufSize := opts.LineBufferSize
	if bufSize <= 0 {
		bufSize = 64 * 1024
	}
	fmt := opts.Format
	if fmt == "" {
		fmt = FormatRaw
	}
	return &Writer{
		bw:     bufio.NewWriterSize(dest, bufSize),
		format: fmt,
	}
}

// WriteLine writes a single log line to the underlying writer.
func (w *Writer) WriteLine(line string) error {
	w.lineNum++
	var err error
	switch w.format {
	case FormatNumbered:
		_, err = fmt.Fprintf(w.bw, "%d\t%s\n", w.lineNum, line)
	default:
		_, err = fmt.Fprintln(w.bw, line)
	}
	return err
}

// Flush flushes any buffered data to the underlying writer.
func (w *Writer) Flush() error {
	return w.bw.Flush()
}

// LinesWritten returns the total number of lines written.
func (w *Writer) LinesWritten() int {
	return w.lineNum
}
