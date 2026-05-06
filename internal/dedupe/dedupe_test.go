package dedupe

import (
	"fmt"
	"testing"
)

func TestZeroWindowAllowsAll(t *testing.T) {
	d := New(0)
	for i := 0; i < 5; i++ {
		if d.IsDuplicate("same line") {
			t.Fatal("zero-window deduper should never report duplicate")
		}
	}
}

func TestFirstOccurrenceNotDuplicate(t *testing.T) {
	d := New(10)
	if d.IsDuplicate("hello") {
		t.Fatal("first occurrence must not be a duplicate")
	}
}

func TestSecondOccurrenceIsDuplicate(t *testing.T) {
	d := New(10)
	d.IsDuplicate("hello")
	if !d.IsDuplicate("hello") {
		t.Fatal("second occurrence within window must be a duplicate")
	}
}

func TestDistinctLinesNotDuplicate(t *testing.T) {
	d := New(10)
	for i := 0; i < 10; i++ {
		line := fmt.Sprintf("line-%d", i)
		if d.IsDuplicate(line) {
			t.Fatalf("unique line %q reported as duplicate", line)
		}
	}
}

func TestWindowEviction(t *testing.T) {
	// Window of 3: after 3 new lines, the first should be evictable.
	d := New(3)
	d.IsDuplicate("a") // slot 0
	d.IsDuplicate("b") // slot 1
	d.IsDuplicate("c") // slot 2 — window full, next insert evicts "a"
	d.IsDuplicate("d") // evicts "a"
	// "a" should no longer be considered a duplicate.
	if d.IsDuplicate("a") {
		t.Fatal("evicted line should not be reported as duplicate")
	}
}

func TestReset(t *testing.T) {
	d := New(5)
	d.IsDuplicate("x")
	d.Reset()
	if d.IsDuplicate("x") {
		t.Fatal("after Reset, previously seen line must not be a duplicate")
	}
}

func TestHashCollisionSafety(t *testing.T) {
	// Ensure two clearly different strings are not conflated.
	d := New(4)
	d.IsDuplicate("foo")
	if d.IsDuplicate("bar") {
		t.Fatal("different lines must not be treated as duplicates")
	}
}
