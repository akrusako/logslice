package coalesce

import (
	"testing"
)

func TestNewInvalidPatternReturnsError(t *testing.T) {
	_, err := New("[", " ")
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestSingleAnchorFlushedAtEnd(t *testing.T) {
	c, _ := New(`^\s+`, " ")
	out, ok := c.Process("anchor line")
	if ok {
		t.Fatalf("expected no output yet, got %q", out)
	}
	flushed, ok := c.Flush()
	if !ok || flushed != "anchor line" {
		t.Fatalf("expected 'anchor line', got %q ok=%v", flushed, ok)
	}
}

func TestContinuationMergedIntoAnchor(t *testing.T) {
	c, _ := New(`^\s+`, " ")
	c.Process("ERROR: something went wrong")
	c.Process("  at foo.go:10")
	c.Process("  at bar.go:20")

	// Flush to get the merged entry.
	result, ok := c.Flush()
	if !ok {
		t.Fatal("expected flushed entry")
	}
	expected := "ERROR: something went wrong at foo.go:10 at bar.go:20"
	if result != expected {
		t.Fatalf("expected %q, got %q", expected, result)
	}
	if c.MergedCount() != 2 {
		t.Fatalf("expected 2 merged, got %d", c.MergedCount())
	}
}

func TestNewAnchorFlushesOldEntry(t *testing.T) {
	c, _ := New(`^\s+`, " ")
	c.Process("first anchor")
	c.Process("  continuation")

	// Second anchor should flush the first.
	flushed, ok := c.Process("second anchor")
	if !ok {
		t.Fatal("expected flush when new anchor arrives")
	}
	if flushed != "first anchor continuation" {
		t.Fatalf("unexpected flushed value: %q", flushed)
	}

	// Flush to get second anchor.
	second, ok := c.Flush()
	if !ok || second != "second anchor" {
		t.Fatalf("expected 'second anchor', got %q ok=%v", second, ok)
	}
}

func TestFlushEmptyReturnsFalse(t *testing.T) {
	c, _ := New(`^\s+`, " ")
	_, ok := c.Flush()
	if ok {
		t.Fatal("expected false on empty flush")
	}
}

func TestNoContinuationsEachLineFlushedByNext(t *testing.T) {
	c, _ := New(`^\s+`, " ")
	lines := []string{"line one", "line two", "line three"}
	var results []string
	for _, l := range lines {
		if out, ok := c.Process(l); ok {
			results = append(results, out)
		}
	}
	if out, ok := c.Flush(); ok {
		results = append(results, out)
	}
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d: %v", len(results), results)
	}
	if c.MergedCount() != 0 {
		t.Fatalf("expected 0 merged, got %d", c.MergedCount())
	}
}
