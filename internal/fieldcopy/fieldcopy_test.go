package fieldcopy

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestNewEmptySrcReturnsError(t *testing.T) {
	_, err := New("", "dst")
	if err == nil {
		t.Fatal("expected error for empty src")
	}
}

func TestNewEmptyDstReturnsError(t *testing.T) {
	_, err := New("src", "")
	if err == nil {
		t.Fatal("expected error for empty dst")
	}
}

func TestNewSameSrcDstReturnsError(t *testing.T) {
	_, err := New("field", "field")
	if err == nil {
		t.Fatal("expected error when src == dst")
	}
}

func TestNewValidParams(t *testing.T) {
	c, err := New("src", "dst")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil copier")
	}
}

func TestCopyKVField(t *testing.T) {
	c, _ := New("level", "severity")
	out := c.Apply("ts=2024-01-01 level=error msg=oops")
	if !strings.Contains(out, "severity=error") {
		t.Errorf("expected severity=error in %q", out)
	}
	if !strings.Contains(out, "level=error") {
		t.Errorf("original field should remain in %q", out)
	}
}

func TestCopyKVMissingFieldUnchanged(t *testing.T) {
	c, _ := New("level", "severity")
	line := "ts=2024-01-01 msg=hello"
	out := c.Apply(line)
	if out != line {
		t.Errorf("expected unchanged line, got %q", out)
	}
}

func TestCopyJSONField(t *testing.T) {
	c, _ := New("level", "severity")
	line := `{"ts":"2024-01-01","level":"warn","msg":"hi"}`
	out := c.Apply(line)
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(out), &m); err != nil {
		t.Fatalf("output not valid JSON: %v", err)
	}
	if m["severity"] != "warn" {
		t.Errorf("expected severity=warn, got %v", m["severity"])
	}
	if m["level"] != "warn" {
		t.Errorf("original level should remain, got %v", m["level"])
	}
}

func TestCopyJSONMissingFieldUnchanged(t *testing.T) {
	c, _ := New("level", "severity")
	line := `{"msg":"hello"}`
	out := c.Apply(line)
	if out != line {
		t.Errorf("expected unchanged line, got %q", out)
	}
}

func TestCopiedCounterIncrements(t *testing.T) {
	c, _ := New("level", "severity")
	c.Apply("level=info msg=ok")
	c.Apply("msg=no-level")
	c.Apply(`{"level":"debug"}`)
	if c.Copied() != 2 {
		t.Errorf("expected 2 copies, got %d", c.Copied())
	}
}

func TestPlainTextUnchanged(t *testing.T) {
	c, _ := New("level", "severity")
	line := "just a plain log message"
	out := c.Apply(line)
	if out != line {
		t.Errorf("expected plain text unchanged, got %q", out)
	}
}
