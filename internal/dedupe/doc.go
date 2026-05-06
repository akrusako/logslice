// Package dedupe implements a sliding-window line deduplicator for log
// processing pipelines.
//
// It is designed to suppress bursts of identical log lines (e.g. repeated
// error messages) without buffering the entire input. Deduplication is
// performed by hashing each line with FNV-64a and maintaining a fixed-size
// ring buffer of recent hashes.
//
// Example usage:
//
//	d := dedupe.New(1000)
//	for _, line := range lines {
//		if !d.IsDuplicate(line) {
//			fmt.Println(line)
//		}
//	}
package dedupe
