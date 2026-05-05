// Package tail implements efficient tail-N-lines extraction for log files.
//
// It uses a pre-built line index (see package lineindex) to seek directly
// to the start offset of the desired line, avoiding a full file scan.
//
// Typical usage:
//
//	r, err := tail.New("/var/log/app.log")
//	if err != nil { ... }
//	lines, err := r.Lines(50)
//	if err != nil { ... }
//	for _, l := range lines {
//		fmt.Println(l)
//	}
//
// The Reader caches the line index so multiple calls to Lines are cheap.
package tail
