package throttle

import (
	"testing"
	"time"
)

func TestDisabledAllowsAll(t *testing.T) {
	th := New(0, 0)
	for i := 0; i < 100; i++ {
		if !th.Allow() {
			t.Fatalf("disabled throttle should allow all lines")
		}
	}
	if th.Dropped() != 0 {
		t.Fatalf("expected 0 dropped, got %d", th.Dropped())
	}
}

func TestWindowCapEnforced(t *testing.T) {
	th := New(time.Second, 3)
	allowed := 0
	for i := 0; i < 10; i++ {
		if th.Allow() {
			allowed++
		}
	}
	if allowed != 3 {
		t.Fatalf("expected 3 allowed, got %d", allowed)
	}
	if th.Dropped() != 7 {
		t.Fatalf("expected 7 dropped, got %d", th.Dropped())
	}
}

func TestTotalCountsAllLines(t *testing.T) {
	th := New(time.Second, 2)
	for i := 0; i < 5; i++ {
		th.Allow()
	}
	if th.Total() != 5 {
		t.Fatalf("expected total 5, got %d", th.Total())
	}
}

func TestWindowEvictionResetsCapacity(t *testing.T) {
	th := New(50*time.Millisecond, 2)

	// Fill the window.
	th.Allow()
	th.Allow()
	if th.Allow() {
		t.Fatal("third call should be throttled")
	}

	// Wait for the window to expire.
	time.Sleep(60 * time.Millisecond)

	// Should be allowed again.
	if !th.Allow() {
		t.Fatal("after window expiry, line should be allowed")
	}
}

func TestZeroMaxLinesAllowsAll(t *testing.T) {
	th := New(time.Second, 0)
	for i := 0; i < 50; i++ {
		if !th.Allow() {
			t.Fatalf("zero maxLines should allow all lines")
		}
	}
}

func TestDroppedCountAccumulates(t *testing.T) {
	th := New(time.Second, 1)
	th.Allow() // allowed
	th.Allow() // dropped
	th.Allow() // dropped
	if th.Dropped() != 2 {
		t.Fatalf("expected 2 dropped, got %d", th.Dropped())
	}
}
