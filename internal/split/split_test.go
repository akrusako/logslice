package split

import (
	"fmt"
	"testing"
)

func TestNewInvalidBuckets(t *testing.T) {
	_, err := New(0, 0)
	if err == nil {
		t.Fatal("expected error for buckets=0")
	}
}

func TestNewBucketOutOfRange(t *testing.T) {
	_, err := New(3, 3)
	if err == nil {
		t.Fatal("expected error for bucket >= buckets")
	}
}

func TestNewValidParams(t *testing.T) {
	s, err := New(4, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Buckets() != 4 {
		t.Errorf("expected 4 buckets, got %d", s.Buckets())
	}
	if s.Bucket() != 2 {
		t.Errorf("expected bucket 2, got %d", s.Bucket())
	}
}

func TestSingleBucketAllowsAll(t *testing.T) {
	s, _ := New(1, 0)
	lines := []string{"alpha", "beta", "gamma", "delta"}
	for _, l := range lines {
		if !s.Allow(l) {
			t.Errorf("single bucket should allow all lines, rejected: %q", l)
		}
	}
	if s.Total() != len(lines) {
		t.Errorf("expected total %d, got %d", len(lines), s.Total())
	}
	if s.Matched() != len(lines) {
		t.Errorf("expected matched %d, got %d", len(lines), s.Matched())
	}
}

func TestBucketsPartitionLines(t *testing.T) {
	const n = 4
	lines := make([]string, 40)
	for i := range lines {
		lines[i] = fmt.Sprintf("log line number %d with some content", i)
	}

	splitters := make([]*Splitter, n)
	for i := range splitters {
		s, err := New(n, uint64(i))
		if err != nil {
			t.Fatalf("New(%d,%d): %v", n, i, err)
		}
		splitters[i] = s
	}

	for _, l := range lines {
		count := 0
		for _, s := range splitters {
			if s.Allow(l) {
				count++
			}
		}
		if count != 1 {
			t.Errorf("line %q matched %d buckets, want exactly 1", l, count)
		}
	}

	totalMatched := 0
	for _, s := range splitters {
		totalMatched += s.Matched()
	}
	if totalMatched != len(lines) {
		t.Errorf("total matched across buckets = %d, want %d", totalMatched, len(lines))
	}
}

func TestCountersIncrement(t *testing.T) {
	s, _ := New(2, 0)
	for i := 0; i < 10; i++ {
		s.Allow(fmt.Sprintf("line-%d", i))
	}
	if s.Total() != 10 {
		t.Errorf("expected total 10, got %d", s.Total())
	}
}
