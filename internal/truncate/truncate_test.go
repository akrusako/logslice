package truncate

import (
	"strings"
	"testing"
)

func TestNoTruncationWhenDisabled(t *testing.T) {
	tr := New(0)
	line := strings.Repeat("a", 200)
	if got := tr.Apply(line); got != line {
		t.Fatalf("expected unchanged line, got %q", got)
	}
	if tr.TruncatedCount() != 0 {
		t.Fatalf("expected 0 truncations, got %d", tr.TruncatedCount())
	}
}

func TestShortLineUnchanged(t *testing.T) {
	tr := New(80)
	line := "short line"
	if got := tr.Apply(line); got != line {
		t.Fatalf("expected %q, got %q", line, got)
	}
	if tr.TruncatedCount() != 0 {
		t.Fatalf("expected 0 truncations")
	}
}

func TestExactLengthUnchanged(t *testing.T) {
	tr := New(10)
	line := strings.Repeat("x", 10)
	if got := tr.Apply(line); got != line {
		t.Fatalf("expected unchanged line at exact limit")
	}
}

func TestLongLineTruncated(t *testing.T) {
	tr := New(20)
	line := strings.Repeat("a", 50)
	got := tr.Apply(line)
	if len(got) != 20 {
		t.Fatalf("expected length 20, got %d", len(got))
	}
	if !strings.HasSuffix(got, "...") {
		t.Fatalf("expected suffix '...', got %q", got)
	}
	if tr.TruncatedCount() != 1 {
		t.Fatalf("expected 1 truncation, got %d", tr.TruncatedCount())
	}
}

func TestTruncatedCountAccumulates(t *testing.T) {
	tr := New(10)
	for i := 0; i < 5; i++ {
		tr.Apply(strings.Repeat("b", 20))
	}
	tr.Apply("short")
	if tr.TruncatedCount() != 5 {
		t.Fatalf("expected 5 truncations, got %d", tr.TruncatedCount())
	}
}

func TestSummaryNoTruncations(t *testing.T) {
	tr := New(100)
	if got := tr.Summary(); got != "no lines truncated" {
		t.Fatalf("unexpected summary: %q", got)
	}
}

func TestSummaryWithTruncations(t *testing.T) {
	tr := New(10)
	tr.Apply(strings.Repeat("z", 30))
	s := tr.Summary()
	if !strings.Contains(s, "1 line(s) truncated") {
		t.Fatalf("unexpected summary: %q", s)
	}
}

func TestNegativeMaxLenDisables(t *testing.T) {
	tr := New(-1)
	line := strings.Repeat("c", 100)
	if got := tr.Apply(line); got != line {
		t.Fatalf("expected unchanged line for negative maxLen")
	}
}
