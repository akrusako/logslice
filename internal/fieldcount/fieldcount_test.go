package fieldcount_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/fieldcount"
)

func TestNewNegativeMinReturnsError(t *testing.T) {
	_, err := fieldcount.New(-1, 5)
	if err == nil {
		t.Fatal("expected error for negative min")
	}
}

func TestNewMaxLessThanMinReturnsError(t *testing.T) {
	_, err := fieldcount.New(5, 3)
	if err == nil {
		t.Fatal("expected error when max < min")
	}
}

func TestNewUnboundedMaxOK(t *testing.T) {
	_, err := fieldcount.New(0, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAllowKVLine(t *testing.T) {
	f, _ := fieldcount.New(2, 4)
	if !f.Allow("level=info msg=hello") {
		t.Error("expected line with 2 fields to be allowed")
	}
}

func TestDropKVLineTooFewFields(t *testing.T) {
	f, _ := fieldcount.New(3, -1)
	if f.Allow("level=info") {
		t.Error("expected line with 1 field to be dropped")
	}
}

func TestDropKVLineTooManyFields(t *testing.T) {
	f, _ := fieldcount.New(1, 2)
	if f.Allow("a=1 b=2 c=3") {
		t.Error("expected line with 3 fields to be dropped")
	}
}

func TestAllowJSONLine(t *testing.T) {
	f, _ := fieldcount.New(1, 5)
	if !f.Allow(`{"level":"warn","msg":"oops"}`) {
		t.Error("expected JSON line with 2 fields to be allowed")
	}
}

func TestCountersIncrement(t *testing.T) {
	f, _ := fieldcount.New(2, 3)
	f.Allow("a=1 b=2")   // allowed
	f.Allow("x=9")       // dropped (1 field)
	f.Allow("p=1 q=2 r=3") // allowed

	if f.Allowed() != 2 {
		t.Errorf("expected 2 allowed, got %d", f.Allowed())
	}
	if f.Dropped() != 1 {
		t.Errorf("expected 1 dropped, got %d", f.Dropped())
	}
}

func TestCountHelper(t *testing.T) {
	n := fieldcount.Count("a=1 b=2 c=3")
	if n != 3 {
		t.Errorf("expected 3, got %d", n)
	}
}

func TestCountEmptyLine(t *testing.T) {
	n := fieldcount.Count("")
	if n != 0 {
		t.Errorf("expected 0, got %d", n)
	}
}
