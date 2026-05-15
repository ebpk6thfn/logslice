// Package progress provides lightweight progress reporting for logslice
// operations. It tracks how many bytes of a log file have been scanned
// and can emit periodic updates to any io.Writer (e.g. stderr).
//
// Typical usage:
//
//	// Obtain file size before slicing.
//	info, _ := os.Stat(path)
//	rep := progress.New(os.Stderr, info.Size(), time.Second)
//	defer rep.Stop(true)
//
//	// Advance as bytes are consumed.
//	rep.Advance(int64(len(line)))
//
// Reporting is safe for concurrent use; Advance may be called from
// multiple goroutines simultaneously.
package progress
