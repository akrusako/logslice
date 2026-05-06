package ratelimit

import (
	"testing"
	"time"
)

func TestNoLimitAllowsAll(t *testing.T) {
	l := New(0, 0)
	for i := 0; i < 1000; i++ {
		if !l.Allow() {
			t.Fatalf("expected Allow()=true at i=%d", i)
		}
	}
	if l.Total() != 1000 {
		t.Fatalf("expected total 1000, got %d", l.Total())
	}
}

func TestMaxTotalCap(t *testing.T) {
	l := New(0, 5)
	allowed := 0
	for i := 0; i < 20; i++ {
		if l.Allow() {
			allowed++
		}
	}
	if allowed != 5 {
		t.Fatalf("expected 5 allowed, got %d", allowed)
	}
}

func TestTotalCounterIncrements(t *testing.T) {
	l := New(0, 0)
	for i := 1; i <= 10; i++ {
		l.Allow()
		if l.Total() != i {
			t.Fatalf("expected total %d, got %d", i, l.Total())
		}
	}
}

func TestRefillAddsTokens(t *testing.T) {
	base := time.Now()
	calls := 0
	times := []time.Time{
		base,
		base.Add(2 * time.Second), // +2 s → +20 tokens at rate 10
	}
	l := New(10, 0)
	l.lastTick = base
	l.bucket = 0
	l.now = func() time.Time {
		t := times[calls]
		if calls < len(times)-1 {
			calls++
		}
		return t
	}
	l.refill()
	if l.bucket != 10 { // capped at maxPerSec
		t.Fatalf("expected bucket=10 after refill, got %d", l.bucket)
	}
}

func TestMaxTotalBeatsRateLimit(t *testing.T) {
	// maxTotal should stop emission even when tokens are available
	l := New(100, 3)
	allowed := 0
	for i := 0; i < 50; i++ {
		if l.Allow() {
			allowed++
		}
	}
	if allowed != 3 {
		t.Fatalf("expected 3, got %d", allowed)
	}
}
