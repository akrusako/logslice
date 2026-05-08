package columnar

import (
	"strings"
	"testing"
)

func TestNoFieldsReturnsOriginal(t *testing.T) {
	c := New(nil, 0, "\t")
	out, ok := c.Format(`level=info msg=hello`)
	if !ok {
		t.Fatal("expected ok=true for no-field passthrough")
	}
	if out != `level=info msg=hello` {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestExtractKVFields(t *testing.T) {
	c := New([]string{"level", "msg"}, 0, "|")
	out, ok := c.Format(`level=info msg=hello ts=2024-01-01`)
	if !ok {
		t.Fatalf("expected all fields found, got ok=false; out=%q", out)
	}
	parts := strings.Split(out, "|")
	if len(parts) != 2 {
		t.Fatalf("expected 2 columns, got %d", len(parts))
	}
	if parts[0] != "info" {
		t.Errorf("col 0: want 'info', got %q", parts[0])
	}
	if parts[1] != "hello" {
		t.Errorf("col 1: want 'hello', got %q", parts[1])
	}
}

func TestExtractJSONFields(t *testing.T) {
	c := New([]string{"level", "msg"}, 0, ",")
	out, ok := c.Format(`{"level":"warn","msg":"oops","ts":"2024-01-01"}`)
	if !ok {
		t.Fatalf("expected all fields found; out=%q", out)
	}
	if out != "warn,oops" {
		t.Errorf("want 'warn,oops', got %q", out)
	}
}

func TestMissingFieldBlankColumn(t *testing.T) {
	c := New([]string{"level", "missing"}, 0, "|")
	out, ok := c.Format(`level=info msg=hello`)
	if ok {
		t.Error("expected ok=false when a field is missing")
	}
	parts := strings.Split(out, "|")
	if parts[1] != "" {
		t.Errorf("missing field should be blank, got %q", parts[1])
	}
}

func TestDroppedCountIncrements(t *testing.T) {
	c := New([]string{"level", "missing"}, 0, "\t")
	c.Format(`level=info`)
	c.Format(`level=warn`)
	if c.DroppedCount() != 2 {
		t.Errorf("want dropped=2, got %d", c.DroppedCount())
	}
}

func TestPaddedColumns(t *testing.T) {
	c := New([]string{"level", "msg"}, 10, "|")
	out, _ := c.Format(`level=info msg=hi`)
	parts := strings.Split(out, "|")
	for _, p := range parts {
		if len(p) != 10 {
			t.Errorf("expected padded width 10, got %d for %q", len(p), p)
		}
	}
}

func TestHeader(t *testing.T) {
	c := New([]string{"level", "msg"}, 0, "|")
	h := c.Header()
	if h != "level|msg" {
		t.Errorf("want 'level|msg', got %q", h)
	}
}

func TestHeaderPadded(t *testing.T) {
	c := New([]string{"level", "msg"}, 8, "|")
	h := c.Header()
	parts := strings.Split(h, "|")
	for _, p := range parts {
		if len(p) != 8 {
			t.Errorf("header column not padded to 8, got %d: %q", len(p), p)
		}
	}
}
