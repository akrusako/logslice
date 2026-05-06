// Package ratelimit provides output throttling for logslice,
// allowing callers to cap the number of lines emitted per second
// and the total number of lines written in a single run.
package ratelimit

import (
	"time"
)

// Limiter enforces an optional lines-per-second rate and an optional
// maximum total line count.
type Limiter struct {
	maxPerSec int
	maxTotal  int
	total     int
	bucket    int
	lastTick  time.Time
	now       func() time.Time // injectable for tests
}

// New creates a Limiter.
// maxPerSec <= 0 disables rate limiting.
// maxTotal  <= 0 disables total-line capping.
func New(maxPerSec, maxTotal int) *Limiter {
	return &Limiter{
		maxPerSec: maxPerSec,
		maxTotal:  maxTotal,
		bucket:    maxPerSec,
		lastTick:  time.Now(),
		now:       time.Now,
	}
}

// Allow reports whether the next line may be emitted.
// When rate-limiting is active it blocks until a token is available.
func (l *Limiter) Allow() bool {
	if l.maxTotal > 0 && l.total >= l.maxTotal {
		return false
	}
	if l.maxPerSec > 0 {
		l.refill()
		for l.bucket <= 0 {
			time.Sleep(time.Millisecond)
			l.refill()
		}
		l.bucket--
	}
	l.total++
	return true
}

// Total returns the number of lines that have been allowed so far.
func (l *Limiter) Total() int { return l.total }

// refill adds tokens proportional to elapsed time since the last call.
func (l *Limiter) refill() {
	now := l.now()
	elapsed := now.Sub(l.lastTick)
	if elapsed <= 0 {
		return
	}
	add := int(elapsed.Seconds() * float64(l.maxPerSec))
	if add > 0 {
		l.bucket += add
		if l.bucket > l.maxPerSec {
			l.bucket = l.maxPerSec
		}
		l.lastTick = now
	}
}
