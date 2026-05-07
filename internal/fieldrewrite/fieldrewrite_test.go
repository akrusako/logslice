package fieldrewrite

import (
	"testing"
)

func TestNewEmptyRulesOK(t *testing.T) {
	rw, err := New(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rw == nil {
		t.Fatal("expected non-nil Rewriter")
	}
}

func TestNewEmptyFieldReturnsError(t *testing.T) {
	_, err := New([]Rule{{Field: "", Old: "a", New: "b"}})
	if err == nil {
		t.Fatal("expected error for empty field name")
	}
}

func TestNoRulesReturnsOriginal(t *testing.T) {
	rw, _ := New(nil)
	line := `{"level":"info","msg":"hello"}`
	if got := rw.Apply(line); got != line {
		t.Fatalf("expected %q, got %q", line, got)
	}
}

func TestRewriteJSONField(t *testing.T) {
	rw, _ := New([]Rule{{Field: "msg", Old: "hello", New: "world"}})
	line := `{"msg":"say hello please"}`
	got := rw.Apply(line)
	if !containsKV(got, `"msg"`, `"say world please"`) {
		t.Fatalf("expected rewritten msg in %q", got)
	}
}

func TestRewriteJSONNoMatchUnchanged(t *testing.T) {
	rw, _ := New([]Rule{{Field: "msg", Old: "absent", New: "x"}})
	line := `{"msg":"hello"}`
	got := rw.Apply(line)
	// value should be unchanged
	if !containsKV(got, `"msg"`, `"hello"`) {
		t.Fatalf("unexpected change: %q", got)
	}
}

func TestRewriteKVField(t *testing.T) {
	rw, _ := New([]Rule{{Field: "env", Old: "prod", New: "staging"}})
	line := "level=info env=prod host=web01"
	got := rw.Apply(line)
	expected := "level=info env=staging host=web01"
	if got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}

func TestRewriteKVNoMatchUnchanged(t *testing.T) {
	rw, _ := New([]Rule{{Field: "env", Old: "dev", New: "staging"}})
	line := "level=info env=prod"
	got := rw.Apply(line)
	if got != line {
		t.Fatalf("expected %q, got %q", line, got)
	}
}

func TestPlainTextUnchanged(t *testing.T) {
	rw, _ := New([]Rule{{Field: "msg", Old: "hello", New: "world"}})
	line := "this is a plain log line with no structure"
	got := rw.Apply(line)
	if got != line {
		t.Fatalf("expected %q, got %q", line, got)
	}
}

func TestMultipleRulesApplied(t *testing.T) {
	rw, _ := New([]Rule{
		{Field: "env", Old: "prod", New: "staging"},
		{Field: "level", Old: "error", New: "warn"},
	})
	line := "level=error env=prod"
	got := rw.Apply(line)
	expected := "level=warn env=staging"
	if got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}

// containsKV is a small helper to check that both key and value substrings
// appear in the output string (order-independent JSON key check).
func containsKV(s, key, val string) bool {
	return contains(s, key) && contains(s, val)
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		(func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		})())
}
