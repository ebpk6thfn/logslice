// Package filter provides optional post-scan filtering for log lines extracted
// by logslice. After the time-bounded window has been identified via binary
// search, individual lines within that window may be further narrowed by:
//
//   - Log level  – retains only lines whose text contains the specified level
//     keyword (e.g. "ERROR", "WARN", "INFO"). Matching is case-insensitive.
//
//   - Regex pattern – retains only lines whose text matches the compiled
//     regular expression.
//
// Both criteria are AND-combined: a line must satisfy every active criterion to
// be included in the output. When no criteria are set, Filter.IsNoop() returns
// true and callers may skip the Match call entirely for performance.
//
// Usage:
//
//	f, err := filter.New(filter.Options{Level: "ERROR", Pattern: `user=\d+`})
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, line := range lines {
//		if f.Match(line) {
//			fmt.Println(line)
//		}
//	}
package filter
