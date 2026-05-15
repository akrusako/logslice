package linenum

import (
	"testing"
)

func TestDefaultFormatStartsAtOne(t *testing.T) {
	inj := New(1, "")
	out := inj.Process("hello world")
	if out != "[1] hello world" {
		t.Fatalf("expected '[1] hello world', got %q", out)
	}
}

func TestCustomFormat(t *testing.T) {
	inj := New(0, "%d: ")
	out := inj.Process("line")
	if out != "0: line" {
		t.Fatalf("expected '0: line', got %q", out)
	}
}

func TestCounterIncrements(t *testing.T) {
	inj := New(1, "")
	for i := 0; i < 5; i++ {
		inj.Process("x")
	}
	if inj.Count() != 5 {
		t.Fatalf("expected count 5, got %d", inj.Count())
	}
}

func TestCurrentAdvances(t *testing.T) {
	inj := New(10, "")
	inj.Process("a")
	inj.Process("b")
	if inj.Current() != 12 {
		t.Fatalf("expected current 12, got %d", inj.Current())
	}
}

func TestEmptyLineSkipped(t *testing.T) {
	inj := New(1, "")
	out := inj.Process("")
	if out != "" {
		t.Fatalf("expected empty string, got %q", out)
	}
	if inj.Count() != 0 {
		t.Fatalf("expected count 0 after empty line, got %d", inj.Count())
	}
	if inj.Current() != 1 {
		t.Fatalf("expected current still 1, got %d", inj.Current())
	}
}

func TestResetRestoresStart(t *testing.T) {
	inj := New(5, "")
	inj.Process("a")
	inj.Process("b")
	inj.Reset()
	if inj.Count() != 0 {
		t.Fatalf("expected count 0 after reset, got %d", inj.Count())
	}
	if inj.Current() != 5 {
		t.Fatalf("expected current 5 after reset, got %d", inj.Current())
	}
	out := inj.Process("z")
	if out != "[5] z" {
		t.Fatalf("expected '[5] z' after reset, got %q", out)
	}
}

func TestSequentialNumbering(t *testing.T) {
	inj := New(1, "%d>")
	lines := []string{"alpha", "beta", "gamma"}
	expected := []string{"1>alpha", "2>beta", "3>gamma"}
	for i, l := range lines {
		out := inj.Process(l)
		if out != expected[i] {
			t.Fatalf("line %d: expected %q, got %q", i, expected[i], out)
		}
	}
}
