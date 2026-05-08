package redact

import (
	"strings"
	"testing"
)

func TestNoPatternsReturnsOriginal(t *testing.T) {
	r, err := New(nil, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out, changed := r.Apply("hello world")
	if changed {
		t.Error("expected no change")
	}
	if out != "hello world" {
		t.Errorf("got %q", out)
	}
}

func TestInvalidPatternReturnsError(t *testing.T) {
	_, err := New([]string{"[invalid"}, "")
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestDefaultMask(t *testing.T) {
	r, _ := New([]string{`\d+`}, "")
	if r.Mask() != "[REDACTED]" {
		t.Errorf("unexpected default mask: %q", r.Mask())
	}
}

func TestCustomMask(t *testing.T) {
	r, _ := New([]string{`\d+`}, "***")
	out, changed := r.Apply("port=8080")
	if !changed {
		t.Error("expected change")
	}
	if !strings.Contains(out, "***") {
		t.Errorf("mask not applied: %q", out)
	}
}

func TestRedactEmailPattern(t *testing.T) {
	r, err := New([]string{`[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}`}, "[EMAIL]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out, changed := r.Apply("user=alice@example.com logged in")
	if !changed {
		t.Error("expected redaction")
	}
	if strings.Contains(out, "alice@example.com") {
		t.Errorf("email not redacted: %q", out)
	}
}

func TestCountIncrementsOnce(t *testing.T) {
	r, _ := New([]string{`\d+`}, "")
	r.Apply("a=1 b=2")
	r.Apply("no digits here")
	r.Apply("c=3")
	if r.Count() != 2 {
		t.Errorf("expected count 2, got %d", r.Count())
	}
}

func TestMultiplePatternsAllApplied(t *testing.T) {
	r, _ := New([]string{`\d+`, `foo`}, "X")
	out, _ := r.Apply("foo=123")
	if strings.Contains(out, "123") || strings.Contains(out, "foo") {
		t.Errorf("not fully redacted: %q", out)
	}
}

func TestNoMatchDoesNotIncrementCount(t *testing.T) {
	r, _ := New([]string{`\d+`}, "")
	r.Apply("no numbers")
	if r.Count() != 0 {
		t.Errorf("expected count 0, got %d", r.Count())
	}
}
