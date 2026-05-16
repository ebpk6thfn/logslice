// Package merge implements multi-source log stream merging.
//
// It provides a heap-based Merger that consumes multiple pre-sorted Entry
// channels and emits a single chronologically ordered stream, as well as a
// Feed helper that turns an io.Reader into an Entry channel by parsing the
// leading timestamp from each line.
//
// Typical usage:
//
//	src1 := merge.Feed(file1, merge.FeedOptions{})
//	src2 := merge.Feed(file2, merge.FeedOptions{})
//
//	m := merge.New(src1, src2)
//	for entry := range m.Merge() {
//		fmt.Printf("%s\n", entry.Line)
//	}
//
// The Merger assumes each individual source is already ordered by timestamp;
// it does not sort within a single source.
package merge
