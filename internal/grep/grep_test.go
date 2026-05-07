package grep

import (
	"testing"
)

func TestNoPatternMatchesAll(t *testing.T) {
	m, err := New(nil, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := []string{"hello world", "foo bar", ""}
	for _, l := range lines {
		if !m.Match(l) {
			t.Errorf("expected match for %q", l)
		}
	}
	if m.Matched() != int64(len(lines)) {
		t.Errorf("matched = %d, want %d", m.Matched(), len(lines))
	}
}

func TestPatternKeepsMatchingLines(t *testing.T) {
	m, err := New([]string{"error"}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !m.Match("an error occurred") {
		t.Error("expected match")
	}
	if m.Match("all good") {
		t.Error("expected no match")
	}
	if m.Matched() != 1 || m.Dropped() != 1 {
		t.Errorf("counters wrong: matched=%d dropped=%d", m.Matched(), m.Dropped())
	}
}

func TestInvertDropsMatchingLines(t *testing.T) {
	m, err := New([]string{"debug"}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Match("debug: verbose info") {
		t.Error("expected inverted match to drop line")
	}
	if !m.Match("error: something failed") {
		t.Error("expected non-matching line to pass through")
	}
}

func TestMultiplePatternsAnyMatch(t *testing.T) {
	m, err := New([]string{"warn", "error", "fatal"}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, l := range []string{"warn: low disk", "error: crash", "fatal: oom"} {
		if !m.Match(l) {
			t.Errorf("expected match for %q", l)
		}
	}
	if m.Match("info: all fine") {
		t.Error("expected no match for info line")
	}
}

func TestInvalidPatternReturnsError(t *testing.T) {
	_, err := New([]string{"[invalid"}, false)
	if err == nil {
		t.Fatal("expected error for invalid regex, got nil")
	}
}

func TestCaseInsensitivePattern(t *testing.T) {
	m, err := New([]string{"(?i)error"}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, l := range []string{"ERROR", "Error", "error"} {
		if !m.Match(l) {
			t.Errorf("expected case-insensitive match for %q", l)
		}
	}
}
