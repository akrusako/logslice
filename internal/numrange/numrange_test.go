package numrange_test

import (
	"testing"

	"github.com/user/logslice/internal/numrange"
)

func TestNewEmptyFieldReturnsError(t *testing.T) {
	_, err := numrange.New("", 0, 10)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewMinGreaterThanMaxReturnsError(t *testing.T) {
	_, err := numrange.New("latency", 100, 1)
	if err == nil {
		t.Fatal("expected error when min > max")
	}
}

func TestNewValidParams(t *testing.T) {
	_, err := numrange.New("latency", 0, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAllowKVInRange(t *testing.T) {
	f, _ := numrange.New("latency", 10, 200)
	if !f.Allow("level=info latency=50 msg=ok") {
		t.Error("expected line in range to be allowed")
	}
}

func TestDropKVBelowMin(t *testing.T) {
	f, _ := numrange.New("latency", 10, 200)
	if f.Allow("level=info latency=5 msg=ok") {
		t.Error("expected line below min to be dropped")
	}
}

func TestDropKVAboveMax(t *testing.T) {
	f, _ := numrange.New("latency", 10, 200)
	if f.Allow("level=info latency=999 msg=ok") {
		t.Error("expected line above max to be dropped")
	}
}

func TestAllowJSONInRange(t *testing.T) {
	f, _ := numrange.New("duration", 0, 5)
	if !f.Allow(`{"level":"info","duration":3.14,"msg":"ok"}`) {
		t.Error("expected JSON line in range to be allowed")
	}
}

func TestDropJSONOutOfRange(t *testing.T) {
	f, _ := numrange.New("duration", 0, 5)
	if f.Allow(`{"level":"info","duration":7.5,"msg":"ok"}`) {
		t.Error("expected JSON line out of range to be dropped")
	}
}

func TestMissingFieldPassesThrough(t *testing.T) {
	f, _ := numrange.New("latency", 0, 100)
	if !f.Allow("level=info msg=no-latency-here") {
		t.Error("expected line without field to pass through")
	}
}

func TestCounters(t *testing.T) {
	f, _ := numrange.New("n", 1, 3)
	f.Allow("n=2")  // pass
	f.Allow("n=5")  // drop
	f.Allow("n=0")  // drop
	f.Allow("x=1")  // pass (missing field)
	if f.Passed() != 2 {
		t.Errorf("want passed=2, got %d", f.Passed())
	}
	if f.Dropped() != 2 {
		t.Errorf("want dropped=2, got %d", f.Dropped())
	}
}
