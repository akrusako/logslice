package linerange

import (
	"testing"
)

func TestParseSingleLine(t *testing.T) {
	r, err := Parse("5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Start != 5 || r.End != 5 {
		t.Errorf("expected 5:5, got %v", r)
	}
}

func TestParseClosedRange(t *testing.T) {
	r, err := Parse("3:10")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Start != 3 || r.End != 10 {
		t.Errorf("expected 3:10, got %v", r)
	}
}

func TestParseOpenEnd(t *testing.T) {
	r, err := Parse("7:")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Start != 7 || r.End != 0 {
		t.Errorf("expected 7:0, got %v", r)
	}
}

func TestParseOpenStart(t *testing.T) {
	r, err := Parse(":20")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Start != 1 || r.End != 20 {
		t.Errorf("expected 1:20, got %v", r)
	}
}

func TestParseInvalidEmpty(t *testing.T) {
	_, err := Parse("")
	if err == nil {
		t.Fatal("expected error for empty string")
	}
}

func TestParseInvalidEndBeforeStart(t *testing.T) {
	_, err := Parse("10:5")
	if err == nil {
		t.Fatal("expected error when end < start")
	}
}

func TestParseInvalidZero(t *testing.T) {
	_, err := Parse("0")
	if err == nil {
		t.Fatal("expected error for line number 0")
	}
}

func TestContainsWithinRange(t *testing.T) {
	r := Range{Start: 3, End: 7}
	for _, line := range []int{3, 5, 7} {
		if !r.Contains(line) {
			t.Errorf("expected line %d to be in range %v", line, r)
		}
	}
}

func TestContainsOutsideRange(t *testing.T) {
	r := Range{Start: 3, End: 7}
	for _, line := range []int{1, 2, 8, 100} {
		if r.Contains(line) {
			t.Errorf("expected line %d to be outside range %v", line, r)
		}
	}
}

func TestContainsUnboundedEnd(t *testing.T) {
	r := Range{Start: 5, End: 0}
	for _, line := range []int{5, 100, 99999} {
		if !r.Contains(line) {
			t.Errorf("expected line %d to be in unbounded range %v", line, r)
		}
	}
	if r.Contains(4) {
		t.Error("expected line 4 to be outside unbounded range starting at 5")
	}
}

func TestStringRepresentation(t *testing.T) {
	if s := (Range{Start: 3, End: 7}).String(); s != "3:7" {
		t.Errorf("unexpected string: %q", s)
	}
	if s := (Range{Start: 5, End: 0}).String(); s != "5:" {
		t.Errorf("unexpected string: %q", s)
	}
}
