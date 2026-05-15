//go:build !windows

package tail

import (
	"os"
	"syscall"
)

// inode returns the inode number for the given FileInfo.
// On non-Unix systems this always returns 0.
func inode(fi os.FileInfo) uint64 {
	if stat, ok := fi.Sys().(*syscall.Stat_t); ok {
		return stat.Ino
	}
	return 0
}
