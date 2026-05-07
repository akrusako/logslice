// Package throttle provides a time-based line emission throttle that limits
// how many lines are forwarded within a sliding time window.
package throttle

import (
	"sync"
	"time"
)

// Throttle tracks line counts within a sliding window and signals when the
// per-window cap has been reached.
type Throttle struct {
	mu       sync.Mutex
	window   time.Duration
	maxLines int
	buckets  []entry
	total    int64
	dropped  int64
}

type entry struct {
	at    time.Time
	count int
}

// New creates a Throttle that allows at most maxLines lines per window.
// If maxLines <= 0 or window <= 0 throttling is disabled and all lines pass.
func New(window time.Duration, maxLines int) *Throttle {
	return &Throttle{
		window:   window,
		maxLines: maxLines,
	}
}

// Allow returns true if the line should be forwarded, false if it should be
// dropped due to the rate cap being exceeded for the current window.
func (t *Throttle) Allow() bool {
	if t.window <= 0 || t.maxLines <= 0 {
		t.total++
		return true
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	t.evict(now)

	windowCount := t.windowCount()
	if windowCount >= t.maxLines {
		t.dropped++
		return false
	}

	t.record(now)
	t.total++
	return true
}

// Dropped returns the total number of lines dropped by this throttle.
func (t *Throttle) Dropped() int64 {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.dropped
}

// Total returns the total number of lines seen (allowed + dropped).
func (t *Throttle) Total() int64 {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.total
}

func (t *Throttle) evict(now time.Time) {
	cutoff := now.Add(-t.window)
	i := 0
	for i < len(t.buckets) && t.buckets[i].at.Before(cutoff) {
		i++
	}
	t.buckets = t.buckets[i:]
}

func (t *Throttle) windowCount() int {
	n := 0
	for _, b := range t.buckets {
		n += b.count
	}
	return n
}

func (t *Throttle) record(now time.Time) {
	if len(t.buckets) > 0 && t.buckets[len(t.buckets)-1].at.Equal(now) {
		t.buckets[len(t.buckets)-1].count++
		return
	}
	t.buckets = append(t.buckets, entry{at: now, count: 1})
}
