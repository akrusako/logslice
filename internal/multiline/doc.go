// Package multiline provides a Merger that reassembles multi-line log entries
// into single logical records.
//
// Many logging frameworks (Java stack traces, Python tracebacks, structured
// multi-line payloads) emit a single logical event across several physical
// lines. multiline.Merger detects the start of a new entry via a configurable
// regular expression and buffers continuation lines until the next entry
// boundary or end-of-stream.
//
// Usage:
//
//	merger, err := multiline.New(`^\d{4}-\d{2}-\d{2}`, " ")
//	for _, line := range lines {
//		if entry, ok := merger.Feed(line); ok {
//			process(entry)
//		}
//	}
//	if entry, ok := merger.Flush(); ok {
//		process(entry)
//	}
package multiline
