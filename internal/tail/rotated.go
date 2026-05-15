package tail

import (
	"context"
	"os"
	"time"
)

// RotatedTailer wraps Tailer and detects log rotation by watching for
// inode / size changes, reopening the file when rotation is detected.
type RotatedTailer struct {
	path         string
	pollInterval time.Duration
	lines        chan Line
}

// NewRotated creates a RotatedTailer that survives log rotation.
func NewRotated(path string, pollInterval time.Duration) *RotatedTailer {
	if pollInterval <= 0 {
		pollInterval = 250 * time.Millisecond
	}
	return &RotatedTailer{
		path:         path,
		pollInterval: pollInterval,
		lines:        make(chan Line, 64),
	}
}

// Lines returns the channel on which new lines are delivered.
func (rt *RotatedTailer) Lines() <-chan Line {
	return rt.lines
}

// Follow tails the file, reopening it if rotation is detected.
func (rt *RotatedTailer) Follow(ctx context.Context) error {
	defer close(rt.lines)

	var lastIno uint64
	var offset int64

	ticker := time.NewTicker(rt.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			info, err := os.Stat(rt.path)
			if err != nil {
				// File may be mid-rotation; retry next tick.
				continue
			}
			curIno := inode(info)
			if curIno != lastIno {
				// Rotation detected or first open — reset offset.
				lastIno = curIno
				offset = 0
			}
			n, err := rt.readFrom(offset)
			offset += n
			if err != nil {
				continue
			}
		}
	}
}

// readFrom reads new content from offset and returns bytes consumed.
func (rt *RotatedTailer) readFrom(offset int64) (int64, error) {
	f, err := os.Open(rt.path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	if _, err := f.Seek(offset, 0); err != nil {
		return 0, err
	}

	var consumed int64
	buf := make([]byte, 0, 512)
	tmp := make([]byte, 512)
	for {
		n, err := f.Read(tmp)
		for i := 0; i < n; i++ {
			if tmp[i] == '\n' {
				rt.lines <- Line{Text: string(buf), Offset: offset + consumed}
				buf = buf[:0]
			} else {
				buf = append(buf, tmp[i])
			}
		}
		consumed += int64(n)
		if err != nil {
			break
		}
	}
	return consumed, nil
}
