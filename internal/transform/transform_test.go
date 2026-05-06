package transform

import (
	"testing"
)

func TestNoOpReturnsOriginal(t *testing.T) {
	tr := New()
	const line = "hello world"
	if got := tr.Apply(line); got != line {
		t.Fatalf("expected %q, got %q", line, got)
	}
}

func TestUpperCase(t *testing.T) {
	tr := New(WithUpper())
	if got := tr.Apply("hello"); got != "HELLO" {
		t.Fatalf("unexpected %q", got)
	}
}

func TestLowerCase(t *testing.T) {
	tr := New(WithLower())
	if got := tr.Apply("HELLO"); got != "hello" {
		t.Fatalf("unexpected %q", got)
	}
}

func TestStripPrefix(t *testing.T) {
	tr := New(WithStripPrefix("INFO "))
	if got := tr.Apply("INFO hello"); got != "hello" {
		t.Fatalf("unexpected %q", got)
	}
}

func TestStripPrefixNoMatch(t *testing.T) {
	tr := New(WithStripPrefix("INFO "))
	const line = "DEBUG hello"
	if got := tr.Apply(line); got != line {
		t.Fatalf("unexpected %q", got)
	}
}

func TestStripSuffix(t *testing.T) {
	tr := New(WithStripSuffix(" END"))
	if got := tr.Apply("message END"); got != "message" {
		t.Fatalf("unexpected %q", got)
	}
}

func TestRenameJSONField(t *testing.T) {
	tr := New(WithRenameFields(map[string]string{"msg": "message"}))
	input := `{"msg":"hello","level":"info"}`
	got := tr.Apply(input)
	if !containsKey(got, "message") {
		t.Fatalf("expected 'message' key in %q", got)
	}
	if containsKey(got, `"msg"`) {
		t.Fatalf("old key 'msg' should be gone in %q", got)
	}
}

func TestRenameKVField(t *testing.T) {
	tr := New(WithRenameFields(map[string]string{"lvl": "level"}))
	got := tr.Apply("lvl=info msg=hello")
	if got != "level=info msg=hello" {
		t.Fatalf("unexpected %q", got)
	}
}

func TestRenameKVNoMatch(t *testing.T) {
	tr := New(WithRenameFields(map[string]string{"x": "y"}))
	const line = "lvl=info msg=hello"
	if got := tr.Apply(line); got != line {
		t.Fatalf("unexpected %q", got)
	}
}

func containsKey(s, key string) bool {
	return len(s) > 0 && (len(key) == 0 || (func() bool {
		for i := 0; i <= len(s)-len(key); i++ {
			if s[i:i+len(key)] == key {
				return true
			}
		}
		return false
	}()))
}
