// Package burst detects and flags bursts of log lines within a sliding
// time window. A burst is defined as more than MaxLines lines arriving
// within WindowDuration. Once a burst is detected the Detector returns
// true so callers can annotate or suppress the excess lines.
package burst

import "time"

// Detector tracks line arrival times and reports burst conditions.
type Detector struct {
	window   time.Duration
	maxLines int
	times    []time.Time
	total    int
	bursts   int
}

// New creates a Detector that fires when more than maxLines lines are
// seen within the given window duration. A maxLines value of 0 disables
// burst detection and IsBurst always returns false.
func New(window time.Duration, maxLines int) *Detector {
	return &Detector{
		window:   window,
		maxLines: maxLines,
		times:    make([]time.Time, 0, 64),
	}
}

// Record registers a new line arrival at the given timestamp and returns
// true when the arrival triggers a burst condition.
func (d *Detector) Record(at time.Time) bool {
	d.total++
	if d.maxLines == 0 {
		return false
	}

	// Evict entries outside the window.
	cutoff := at.Add(-d.window)
	i := 0
	for i < len(d.times) && d.times[i].Before(cutoff) {
		i++
	}
	d.times = append(d.times[:0], d.times[i:]...)
	d.times = append(d.times, at)

	if len(d.times) > d.maxLines {
		d.bursts++
		return true
	}
	return false
}

// Total returns the total number of lines recorded.
func (d *Detector) Total() int { return d.total }

// Bursts returns the number of lines that triggered a burst condition.
func (d *Detector) Bursts() int { return d.bursts }

// WindowSize returns the configured sliding window duration.
func (d *Detector) WindowSize() time.Duration { return d.window }
