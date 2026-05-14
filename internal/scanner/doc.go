// Package scanner provides line-oriented log scanning with timestamp parsing
// and byte-offset tracking.
//
// The primary types are:
//
//   - Scanner: reads an io.Reader line by line, parsing timestamps and recording
//     byte offsets for each line. Designed for sequential forward passes.
//
//   - Seeker: wraps an io.ReadSeeker to locate the start and end byte offsets
//     of a time-bounded window within a log file. Used by the slicer to avoid
//     loading the entire file into memory.
//
// Typical usage:
//
//	f, _ := os.Open("app.log")
//	defer f.Close()
//
//	// Sequential scan
//	s := scanner.New(f, time.RFC3339)
//	for s.Scan() {
//		line := s.Line()
//		fmt.Printf("ts=%v offset=%d\n", line.Timestamp, line.Offset)
//	}
//
//	// Find time window offsets
//	sk := scanner.NewSeeker(f, time.RFC3339)
//	startOff, _ := sk.FindStart(0, from)
//	endOff, _   := sk.FindEnd(startOff, to)
package scanner
