// Package tail provides functionality to follow a log file and emit
// new lines as they are appended, similar to `tail -f`.
package tail

import (
	"bufio"
	"context"
	"io"
	"os"
	"time"
)

// Line represents a single line emitted by the tailer.
type Line struct {
	Text   string
	Offset int64
}

// Tailer follows a file and sends new lines to a channel.
type Tailer struct {
	path     string
	pollInterval time.Duration
	lines    chan Line
}

// New creates a new Tailer for the given file path.
// pollInterval controls how often the file is checked for new content.
func New(path string, pollInterval time.Duration) *Tailer {
	if pollInterval <= 0 {
		pollInterval = 250 * time.Millisecond
	}
	return &Tailer{
		path:         path,
		pollInterval: pollInterval,
		lines:        make(chan Line, 64),
	}
}

// Lines returns the channel on which new lines are delivered.
func (t *Tailer) Lines() <-chan Line {
	return t.lines
}

// Follow starts tailing the file from its current end.
// It blocks until ctx is cancelled, then closes the Lines channel.
func (t *Tailer) Follow(ctx context.Context) error {
	f, err := os.Open(t.path)
	if err != nil {
		return err
	}
	defer f.Close()

	// Seek to end so we only see new content.
	offset, err := f.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}

	reader := bufio.NewReader(f)
	ticker := time.NewTicker(t.pollInterval)
	defer ticker.Stop()
	defer close(t.lines)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			for {
				line, err := reader.ReadString('\n')
				if len(line) > 0 {
					// Strip trailing newline for consistency.
					text := line
					if len(text) > 0 && text[len(text)-1] == '\n' {
						text = text[:len(text)-1]
					}
					t.lines <- Line{Text: text, Offset: offset}
					offset += int64(len(line))
				}
				if err == io.EOF {
					break
				}
				if err != nil {
					return err
				}
			}
		}
	}
}
