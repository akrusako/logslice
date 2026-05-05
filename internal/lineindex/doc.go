// Package lineindex builds and queries a byte-offset index over a log file,
// allowing O(log n) seeks to any line number or byte position without
// re-scanning the entire file.
//
// Typical usage:
//
//	f, _ := os.Open("app.log")
//	defer f.Close()
//
//	idx, err := lineindex.Build(f)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Jump directly to line 5000
//	offset := idx.OffsetForLine(5000)
//	f.Seek(offset, io.SeekStart)
//
// The index is entirely in-memory and is rebuilt on each invocation of Build.
// For very large files (>1 GB) consider streaming approaches instead.
package lineindex
