package lineindex_test

import (
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/lineindex"
)

const sampleLog = `2024-01-01T00:00:01Z level=info msg="start"
2024-01-01T00:00:02Z level=debug msg="tick"
2024-01-01T00:00:03Z level=warn msg="slow"
2024-01-01T00:00:04Z level=error msg="fail"
`

func buildIndex(t *testing.T) lineindex.Index {
	t.Helper()
	idx, err := lineindex.Build(strings.NewReader(sampleLog))
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	return idx
}

func TestBuildLen(t *testing.T) {
	idx := buildIndex(t)
	if got := idx.Len(); got != 4 {
		t.Errorf("Len() = %d, want 4", got)
	}
}

func TestOffsetForLineFirstLine(t *testing.T) {
	idx := buildIndex(t)
	if off := idx.OffsetForLine(1); off != 0 {
		t.Errorf("OffsetForLine(1) = %d, want 0", off)
	}
}

func TestOffsetForLineOutOfRange(t *testing.T) {
	idx := buildIndex(t)
	if off := idx.OffsetForLine(0); off != -1 {
		t.Errorf("OffsetForLine(0) = %d, want -1", off)
	}
	if off := idx.OffsetForLine(99); off != -1 {
		t.Errorf("OffsetForLine(99) = %d, want -1", off)
	}
}

func TestOffsetForLineMonotonic(t *testing.T) {
	idx := buildIndex(t)
	prev := int64(-1)
	for i := 1; i <= idx.Len(); i++ {
		off := idx.OffsetForLine(i)
		if off <= prev {
			t.Errorf("line %d: offset %d not greater than previous %d", i, off, prev)
		}
		prev = off
	}
}

func TestLineForOffset(t *testing.T) {
	idx := buildIndex(t)
	// offset 0 should be line 1
	if line := idx.LineForOffset(0); line != 1 {
		t.Errorf("LineForOffset(0) = %d, want 1", line)
	}
	// offset exactly at line 2 start
	off2 := idx.OffsetForLine(2)
	if line := idx.LineForOffset(off2); line != 2 {
		t.Errorf("LineForOffset(%d) = %d, want 2", off2, line)
	}
}

func TestRoundTrip(t *testing.T) {
	idx := buildIndex(t)
	for want := 1; want <= idx.Len(); want++ {
		off := idx.OffsetForLine(want)
		got := idx.LineForOffset(off)
		if got != want {
			t.Errorf("round-trip line %d: got %d", want, got)
		}
	}
}
