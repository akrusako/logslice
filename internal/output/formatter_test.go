package output

import (
	"bytes"
	"strings"
	"testing"
)

func TestFormatRaw(t *testing.T) {
	var buf bytes.Buffer
	f := New(&buf, FormatRaw)

	line := `{"level":"info","msg":"hello"}`
	if err := f.WriteLine(line); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := strings.TrimSpace(buf.String())
	if got != line {
		t.Errorf("expected %q, got %q", line, got)
	}
}

func TestFormatJSON(t *testing.T) {
	var buf bytes.Buffer
	f := New(&buf, FormatJSON)

	line := `{"level":"info","msg":"hello"}`
	if err := f.WriteLine(line); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "\n  ") {
		t.Errorf("expected indented JSON, got: %s", out)
	}
}

func TestFormatJSONFallback(t *testing.T) {
	var buf bytes.Buffer
	f := New(&buf, FormatJSON)

	line := "not json at all"
	if err := f.WriteLine(line); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := strings.TrimSpace(buf.String())
	if got != line {
		t.Errorf("expected raw fallback %q, got %q", line, got)
	}
}

func TestFormatTSV(t *testing.T) {
	var buf bytes.Buffer
	f := New(&buf, FormatTSV)

	line := `{"level":"warn","code":42}`
	if err := f.WriteLine(line); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := strings.TrimSpace(buf.String())
	if !strings.Contains(out, "\t") {
		t.Errorf("expected tab-separated output, got: %s", out)
	}
}

func TestFormatTSVFallback(t *testing.T) {
	var buf bytes.Buffer
	f := New(&buf, FormatTSV)

	line := "plain text line"
	if err := f.WriteLine(line); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := strings.TrimSpace(buf.String())
	if got != line {
		t.Errorf("expected raw fallback %q, got %q", line, got)
	}
}
