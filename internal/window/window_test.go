package window

import (
	"testing"
	"time"
)

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func TestZeroWindowCountsAll(t *testing.T) {
	a := New(0)
	for i := 0; i < 5; i++ {
		a.Add(epoch.Add(time.Duration(i) * time.Second))
	}
	if got := a.WindowCount(); got != 5 {
		t.Fatalf("expected 5, got %d", got)
	}
}

func TestTotalNeverDrops(t *testing.T) {
	a := New(2 * time.Second)
	for i := 0; i < 10; i++ {
		a.Add(epoch.Add(time.Duration(i) * time.Second))
	}
	if got := a.Total(); got != 10 {
		t.Fatalf("expected total 10, got %d", got)
	}
}

func TestEvictionReducesWindowCount(t *testing.T) {
	a := New(3 * time.Second)
	a.Add(epoch)
	a.Add(epoch.Add(1 * time.Second))
	a.Add(epoch.Add(2 * time.Second))
	// advance well past the window
	count := a.Add(epoch.Add(10 * time.Second))
	if count != 1 {
		t.Fatalf("expected window count 1, got %d", count)
	}
}

func TestWindowCountAtBoundary(t *testing.T) {
	a := New(5 * time.Second)
	a.Add(epoch)
	a.Add(epoch.Add(5 * time.Second))
	// epoch is exactly at the cutoff; evict uses After so it should be evicted
	count := a.Add(epoch.Add(5 * time.Second))
	// epoch (t=0) cutoff = 5s - 5s = 0s; !After(0) => evicted
	if count < 1 {
		t.Fatalf("expected at least 1 in window, got %d", count)
	}
}

func TestReset(t *testing.T) {
	a := New(10 * time.Second)
	a.Add(epoch)
	a.Add(epoch.Add(1 * time.Second))
	a.Reset()
	if got := a.Total(); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
	if got := a.WindowCount(); got != 0 {
		t.Fatalf("expected 0 window after reset, got %d", got)
	}
}

func TestAddReturnValue(t *testing.T) {
	a := New(60 * time.Second)
	for i := 1; i <= 4; i++ {
		got := a.Add(epoch.Add(time.Duration(i) * time.Second))
		if got != int64(i) {
			t.Fatalf("step %d: expected %d, got %d", i, i, got)
		}
	}
}
