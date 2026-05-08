package head

import (
	"testing"
)

func TestZeroMaxAllowsAll(t *testing.T) {
	p := New(0)
	lines := []string{"a", "b", "c", "d", "e"}
	for _, l := range lines {
		out, ok := p.Process(l)
		if !ok || out != l {
			t.Fatalf("expected line %q to pass through, got ok=%v out=%q", l, ok, out)
		}
	}
	if p.Done() {
		t.Fatal("Done() should be false when max <= 0")
	}
}

func TestNegativeMaxAllowsAll(t *testing.T) {
	p := New(-5)
	_, ok := p.Process("hello")
	if !ok {
		t.Fatal("negative max should allow all lines")
	}
}

func TestMaxOneLimitsToFirstLine(t *testing.T) {
	p := New(1)
	out, ok := p.Process("first")
	if !ok || out != "first" {
		t.Fatalf("expected first line to pass, got ok=%v out=%q", ok, out)
	}
	_, ok = p.Process("second")
	if ok {
		t.Fatal("second line should be dropped")
	}
}

func TestDoneAfterLimit(t *testing.T) {
	p := New(3)
	for i := 0; i < 3; i++ {
		p.Process("line")
	}
	if !p.Done() {
		t.Fatal("Done() should be true after limit is reached")
	}
}

func TestNotDoneBeforeLimit(t *testing.T) {
	p := New(5)
	p.Process("a")
	p.Process("b")
	if p.Done() {
		t.Fatal("Done() should be false before limit is reached")
	}
}

func TestDropCountAccumulates(t *testing.T) {
	p := New(2)
	p.Process("a")
	p.Process("b")
	p.Process("c")
	p.Process("d")
	if p.dropped != 2 {
		t.Fatalf("expected 2 dropped, got %d", p.dropped)
	}
}

func TestStatsString(t *testing.T) {
	p := New(2)
	p.Process("x")
	p.Process("y")
	p.Process("z")
	s := p.Stats()
	if s != "head: kept 2, dropped 1" {
		t.Fatalf("unexpected stats: %q", s)
	}
}
