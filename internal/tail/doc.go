// Package tail provides utilities for following log files in real time,
// including support for log rotation detection.
//
// # Tailer
//
// The basic [Tailer] watches a single file and emits new lines as they are
// appended. It polls the file at a configurable interval and sends each new
// line over a channel.
//
//	tr, err := tail.New("/var/log/app.log", 100*time.Millisecond)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer tr.Stop()
//	for line := range tr.Lines() {
//		fmt.Println(line)
//	}
//
// # RotatedTailer
//
// [NewRotated] extends the basic tailer with inode-based rotation detection.
// When the file at the watched path is replaced (e.g. by logrotate), the
// tailer automatically reopens the new file and continues streaming lines
// without missing entries.
//
//	rt, err := tail.NewRotated("/var/log/app.log", 100*time.Millisecond)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer rt.Stop()
//	for line := range rt.Lines() {
//		fmt.Println(line)
//	}
package tail
