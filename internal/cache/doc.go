// Package cache implements a lightweight, thread-safe in-memory cache for
// log-file offset indexes.
//
// Building an index requires a full sequential scan of a log file, which can
// be expensive for large files.  The cache avoids redundant scans by storing
// the index alongside the file's modification time and byte size.  A cached
// entry is considered valid only when both the modification time and size
// still match the values recorded at build time; any change to either field
// causes the lookup to return a miss so the caller rebuilds the index.
//
// The cache is safe for concurrent use by multiple goroutines.
package cache
