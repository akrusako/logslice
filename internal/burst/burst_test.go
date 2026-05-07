package burst

import (
	"testing"
	"time"
)

func TestDisabledNeverBursts(t *testing.T) {
	d := New(time.Second, 0)
	now := time.Now()
	for i := 0; i < 100; i++ {
		if d.Record(now) {
			t.Fatal("disabled detector should never report burst")
		}
	}
}

func TestBelowThresholdNoBurst(t *testing.T) {
	d := New(time.Second, 5)
	now := time.Now()
	for i := 0; i < 5; i++ {
		if d.Record(now) {
			t.Fatalf("line %d should not trigger burst (threshold 5)", i+1)
		}
	}
}

func TestExceedThresholdTriggersBurst(t *testing.T) {
	d := New(time.Second, 3)
	now := time.Now()
	for i := 0; i < 3; i++ {
		d.Record(now)
	}
	if !d.Record(now) {
		t.Fatal("4th line within window should trigger burst")
	}
}

func TestWindowEvictionResetsBurst(t *testing.T) {
	d := New(500*time.Millisecond, 3)
	base := time.Now()

	// Fill window to threshold.
	for i := 0; i < 3; i++ {
		d.Record(base)
	}

	// Advance past the window; old entries should be evicted.
	later := base.Add(600 * time.Millisecond)
	if d.Record(later) {
		t.Fatal("after window eviction the count should reset; no burst expected")
	}
}

func TestBurstCounterAccumulates(t *testing.T) {
	d := New(time.Second, 2)
	now := time.Now()
	d.Record(now)
	d.Record(now)
	d.Record(now) // burst 1
	d.Record(now) // burst 2

	if got := d.Bursts(); got != 2 {
		t.Fatalf("expected 2 bursts, got %d", got)
	}
}

func TestTotalCounterIncrements(t *testing.T) {
	d := New(time.Second, 10)
	now := time.Now()
	for i := 0; i < 7; i++ {
		d.Record(now)
	}
	if d.Total() != 7 {
		t.Fatalf("expected total 7, got %d", d.Total())
	}
}

func TestWindowSizeReturnsConfigured(t *testing.T) {
	win := 2 * time.Second
	d := New(win, 5)
	if d.WindowSize() != win {
		t.Fatalf("expected window %v, got %v", win, d.WindowSize())
	}
}
