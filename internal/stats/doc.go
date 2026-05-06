// Package stats provides lightweight counters and reporting for logslice
// processing runs.
//
// Usage:
//
//	s := stats.New()
//	for _, line := range lines {
//		matched := filter.Match(line)
//		s.RecordLine(matched, len(line))
//	}
//	s.Finish()
//	s.Print(os.Stderr)
//
// Stats is not safe for concurrent use; callers should synchronise externally
// if lines are processed from multiple goroutines.
package stats
