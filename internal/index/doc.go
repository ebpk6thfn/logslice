// Package index builds and queries a lightweight in-memory offset index for
// log files. Rather than scanning every byte when seeking to a time window,
// callers can build an Index once and then resolve start/end byte offsets in
// O(n) time over the (much smaller) index instead of over the raw file.
//
// Typical usage:
//
//	f, _ := os.Open("app.log")
//	defer f.Close()
//
//	b := index.NewBuilder("") // auto-detect timestamp format
//	idx, err := b.Build(f)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	startOffset := idx.FindStart(from)
//	endOffset   := idx.FindEnd(to)   // -1 means read to EOF
package index
