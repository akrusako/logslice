package multiline

import (
	"testing"
)

func TestNewEmptyPatternReturnsError(t *testing.T) {
	_, err := New("", " ")
	if err == nil {
		t.Fatal("expected error for empty pattern")
	}
}

func TestNewInvalidPatternReturnsError(t *testing.T) {
	_, err := New("[", " ")
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestNewValidParams(t *testing.T) {
	m, err := New(`^\d{4}-`, " ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m == nil {
		t.Fatal("expected non-nil Merger")
	}
}

func TestSingleEntryNoContination(t *testing.T) {
	m, _ := New(`^\d{4}-`, " ")
	entry, ok := m.Feed("2024-01-01 hello")
	if ok || entry != "" {
		t.Fatal("first line should buffer, not flush")
	}
	out, ok2 := m.Flush()
	if !ok2 || out != "2024-01-01 hello" {
		t.Fatalf("expected flushed entry, got %q", out)
	}
}

func TestContinuationMerged(t *testing.T) {
	m, _ := New(`^\d{4}-`, " ")
	m.Feed("2024-01-01 start")
	m.Feed("  continuation one")
	m.Feed("  continuation two")
	out, ok := m.Flush()
	if !ok {
		t.Fatal("expected flush")
	}
	want := "2024-01-01 start   continuation one   continuation two"
	if out != want {
		t.Fatalf("got %q, want %q", out, want)
	}
	if m.Merged() != 2 {
		t.Fatalf("expected 2 merged, got %d", m.Merged())
	}
}

func TestNewStartLineFlushesPrevious(t *testing.T) {
	m, _ := New(`^\d{4}-`, " ")
	m.Feed("2024-01-01 first")
	m.Feed("  cont")
	entry, ok := m.Feed("2024-01-02 second")
	if !ok {
		t.Fatal("expected flush on new start line")
	}
	if entry != "2024-01-01 first   cont" {
		t.Fatalf("unexpected entry: %q", entry)
	}
	out, _ := m.Flush()
	if out != "2024-01-02 second" {
		t.Fatalf("unexpected final entry: %q", out)
	}
}

func TestFlushEmptyReturnsFalse(t *testing.T) {
	m, _ := New(`^\d{4}-`, " ")
	_, ok := m.Flush()
	if ok {
		t.Fatal("flush on empty buffer should return false")
	}
}

func TestFlushedCounter(t *testing.T) {
	m, _ := New(`^\d{4}-`, " ")
	m.Feed("2024-01-01 a")
	m.Feed("2024-01-02 b")
	m.Flush()
	if m.Flushed() != 2 {
		t.Fatalf("expected 2 flushed, got %d", m.Flushed())
	}
}
