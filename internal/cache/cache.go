// Package cache provides a simple in-memory index cache keyed by file path
// and modification time, so repeated slices of the same file reuse the
// previously built offset index without re-scanning from disk.
package cache

import (
	"sync"
	"time"
)

// Entry holds a cached index along with the file metadata used to validate it.
type Entry struct {
	// ModTime is the file modification time at the moment the index was built.
	ModTime time.Time
	// Size is the file size in bytes at the moment the index was built.
	Size int64
	// Offsets contains the byte offsets stored in the index.
	Offsets []int64
	// Timestamps contains the parsed timestamps parallel to Offsets.
	Timestamps []time.Time
}

// Cache is a thread-safe in-memory store for index entries.
type Cache struct {
	mu      sync.RWMutex
	entries map[string]Entry
}

// New returns an initialised, empty Cache.
func New() *Cache {
	return &Cache{
		entries: make(map[string]Entry),
	}
}

// Get returns the cached Entry for path and a boolean indicating whether the
// entry was found and is still valid for the given modTime and size.
func (c *Cache) Get(path string, modTime time.Time, size int64) (Entry, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.entries[path]
	if !ok {
		return Entry{}, false
	}
	if !e.ModTime.Equal(modTime) || e.Size != size {
		return Entry{}, false
	}
	return e, true
}

// Put stores an Entry for path, overwriting any previous entry.
func (c *Cache) Put(path string, entry Entry) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[path] = entry
}

// Invalidate removes the entry for path if it exists.
func (c *Cache) Invalidate(path string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, path)
}

// Len returns the number of entries currently held in the cache.
func (c *Cache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.entries)
}
