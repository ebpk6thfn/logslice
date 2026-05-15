// Package stats provides lightweight metrics collection for logslice operations.
//
// A Collector is created at the start of a slice run and updated as lines are
// scanned, matched, filtered, and written. When the run completes, Finish marks
// the end time so that Elapsed returns an accurate wall-clock duration.
//
// Typical usage:
//
//	col := stats.New()
//	defer col.Finish()
//
//	// ... for each line processed:
//	col.RecordLine(matched, filtered, int64(len(line)))
//	col.RecordWrite(int64(n))
//
//	// Print a summary:
//	col.Print(os.Stderr)
package stats
