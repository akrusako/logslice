// Package window provides a sliding-window line counter for log streams.
//
// An Aggregator tracks how many log lines have been observed within a
// configurable time window. Buckets older than the window size are evicted
// whenever a new line is added, keeping memory usage proportional to the
// number of lines in the active window rather than the total log volume.
//
// A zero window size disables eviction; all lines are retained and counted.
//
// Typical use:
//
//	agg := window.New(30 * time.Second)
//	count := agg.Add(lineTimestamp)
//	if count > threshold {
//		// spike detected
//	}
package window
