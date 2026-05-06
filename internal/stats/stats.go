// Package stats collects and reports processing statistics for a logslice run.
package stats

import (
	"fmt"
	"io"
	"time"
)

// Stats holds counters accumulated during log processing.
type Stats struct {
	LinesScanned  int
	LinesMatched  int
	LinesDropped  int
	BytesRead     int64
	StartTime     time.Time
	EndTime       time.Time
}

// New returns a new Stats instance with StartTime set to now.
func New() *Stats {
	return &Stats{StartTime: time.Now()}
}

// RecordLine updates counters for a single scanned line.
func (s *Stats) RecordLine(matched bool, byteLen int) {
	s.LinesScanned++
	s.BytesRead += int64(byteLen)
	if matched {
		s.LinesMatched++
	} else {
		s.LinesDropped++
	}
}

// Finish marks the end time of processing.
func (s *Stats) Finish() {
	s.EndTime = time.Now()
}

// Elapsed returns the duration between Start and Finish.
func (s *Stats) Elapsed() time.Duration {
	if s.EndTime.IsZero() {
		return time.Since(s.StartTime)
	}
	return s.EndTime.Sub(s.StartTime)
}

// MatchRate returns the fraction of scanned lines that matched, or 0 if none scanned.
func (s *Stats) MatchRate() float64 {
	if s.LinesScanned == 0 {
		return 0
	}
	return float64(s.LinesMatched) / float64(s.LinesScanned)
}

// Print writes a human-readable summary to w.
func (s *Stats) Print(w io.Writer) {
	fmt.Fprintf(w, "scanned=%d matched=%d dropped=%d bytes=%d elapsed=%s match_rate=%.2f%%\n",
		s.LinesScanned,
		s.LinesMatched,
		s.LinesDropped,
		s.BytesRead,
		s.Elapsed().Round(time.Millisecond),
		s.MatchRate()*100,
	)
}
