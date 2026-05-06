package highlight

import (
	"strings"
	"testing"
)

func TestDisabledReturnsOriginal(t *testing.T) {
	h := New("", false)
	line := "error: disk full"
	if got := h.Term(line, "error"); got != line {
		t.Errorf("expected unchanged line, got %q", got)
	}
}

func TestTermWrapsMatch(t *testing.T) {
	h := New(Yellow, true)
	line := "error: disk full"
	got := h.Term(line, "error")
	if !strings.Contains(got, Yellow) {
		t.Error("expected ANSI colour code in output")
	}
	if !strings.Contains(got, "error") {
		t.Error("expected original term to remain in output")
	}
	if !strings.Contains(got, Reset) {
		t.Error("expected reset code in output")
	}
}

func TestTermNoMatchUnchanged(t *testing.T) {
	h := New(Yellow, true)
	line := "everything is fine"
	got := h.Term(line, "error")
	if got != line {
		t.Errorf("expected unchanged line, got %q", got)
	}
}

func TestTermsMultiple(t *testing.T) {
	h := New(Cyan, true)
	line := "warn: retrying connection to host"
	got := h.Terms(line, []string{"warn", "host"})
	if !strings.Contains(got, "warn") {
		t.Error("expected 'warn' in output")
	}
	if !strings.Contains(got, "host") {
		t.Error("expected 'host' in output")
	}
	if strings.Count(got, Cyan) < 2 {
		t.Error("expected at least two colour code occurrences")
	}
}

func TestTermsEmptySlice(t *testing.T) {
	h := New("", true)
	line := "no terms to highlight"
	if got := h.Terms(line, nil); got != line {
		t.Errorf("expected unchanged line, got %q", got)
	}
}

func TestDefaultColourFallback(t *testing.T) {
	h := New("", true)
	if h.colour != Yellow {
		t.Errorf("expected default colour %q, got %q", Yellow, h.colour)
	}
}

func TestFieldHighlight(t *testing.T) {
	h := New(Green, true)
	line := `level="error" msg="disk full"`
	got := h.Field(line, "error")
	if !strings.Contains(got, Green) {
		t.Error("expected green colour code for field value")
	}
}
