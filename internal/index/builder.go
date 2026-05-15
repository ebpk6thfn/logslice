package index

import (
	"os"
	"time"

	"github.com/yourorg/logslice/internal/cache"
)

// CachedBuild returns an Index for the file at path, reusing a previously
// built index from c when the file has not changed since the last build.
//
// On a cache miss the index is built from scratch via Build, stored in c, and
// then returned.  Stat errors and build errors are returned to the caller
// unchanged.
func CachedBuild(path string, format string, c *cache.Cache) (*Index, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if entry, ok := c.Get(path, info.ModTime(), info.Size()); ok {
		return fromCacheEntry(entry), nil
	}

	idx, err := Build(path, format)
	if err != nil {
		return nil, err
	}

	c.Put(path, toCacheEntry(info.ModTime(), info.Size(), idx))
	return idx, nil
}

// toCacheEntry converts an Index into a cache.Entry for storage.
func toCacheEntry(modTime time.Time, size int64, idx *Index) cache.Entry {
	offsets := make([]int64, len(idx.entries))
	timestamps := make([]time.Time, len(idx.entries))
	for i, e := range idx.entries {
		offsets[i] = e.Offset
		timestamps[i] = e.Timestamp
	}
	return cache.Entry{
		ModTime:    modTime,
		Size:       size,
		Offsets:    offsets,
		Timestamps: timestamps,
	}
}

// fromCacheEntry reconstructs an Index from a cache.Entry.
func fromCacheEntry(e cache.Entry) *Index {
	entries := make([]Entry, len(e.Offsets))
	for i := range e.Offsets {
		entries[i] = Entry{
			Offset:    e.Offsets[i],
			Timestamp: e.Timestamps[i],
		}
	}
	return &Index{entries: entries}
}
