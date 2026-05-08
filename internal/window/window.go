package window

import (
	"time"
)

// Aggregator holds counts of lines seen within a sliding time window.
type Aggregator struct {
	size    time.Duration
	buckets []bucket
	total   int64
}

type bucket struct {
	at    time.Time
	count int64
}

// New returns an Aggregator with the given window size.
// A zero or negative size disables windowing (all lines are counted).
func New(size time.Duration) *Aggregator {
	return &Aggregator{size: size}
}

// Add records a line at the given timestamp and returns the current
// count of lines within the window.
func (a *Aggregator) Add(t time.Time) int64 {
	if a.size > 0 {
		a.evict(t)
	}
	a.buckets = append(a.buckets, bucket{at: t, count: 1})
	a.total++
	return a.windowCount()
}

// WindowCount returns the number of lines currently within the window.
func (a *Aggregator) WindowCount() int64 {
	return a.windowCount()
}

// Total returns the total number of lines ever recorded.
func (a *Aggregator) Total() int64 {
	return a.total
}

// Reset clears all recorded data.
func (a *Aggregator) Reset() {
	a.buckets = a.buckets[:0]
	a.total = 0
}

func (a *Aggregator) evict(now time.Time) {
	cutoff := now.Add(-a.size)
	i := 0
	for i < len(a.buckets) && !a.buckets[i].at.After(cutoff) {
		i++
	}
	a.buckets = a.buckets[i:]
}

func (a *Aggregator) windowCount() int64 {
	var n int64
	for _, b := range a.buckets {
		n += b.count
	}
	return n
}
