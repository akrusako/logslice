package mask

import (
	"strings"
	"testing"
)

func TestNoFieldsReturnsOriginal(t *testing.T) {
	m := New(nil, "")
	line := `{"password":"secret","user":"alice"}`
	if got := m.Mask(line); got != line {
		t.Errorf("expected original line, got %q", got)
	}
}

func TestMaskJSONField(t *testing.T) {
	m := New([]string{"password"}, "***")
	line := `{"user":"alice","password":"s3cr3t"}`
	out := m.Mask(line)
	if strings.Contains(out, "s3cr3t") {
		t.Errorf("sensitive value should be redacted, got %q", out)
	}
	if !strings.Contains(out, "***") {
		t.Errorf("expected placeholder in output, got %q", out)
	}
	if !strings.Contains(out, "alice") {
		t.Errorf("non-sensitive field should be preserved, got %q", out)
	}
}

func TestMaskJSONMultipleFields(t *testing.T) {
	m := New([]string{"password", "token"}, "[REDACTED]")
	line := `{"user":"bob","password":"hunter2","token":"abc123"}`
	out := m.Mask(line)
	if strings.Contains(out, "hunter2") || strings.Contains(out, "abc123") {
		t.Errorf("both sensitive values should be redacted, got %q", out)
	}
}

func TestMaskKVField(t *testing.T) {
	m := New([]string{"password"}, "***")
	line := `level=info user=alice password=secret msg=login`
	out := m.Mask(line)
	if strings.Contains(out, "secret") {
		t.Errorf("sensitive value should be redacted in KV line, got %q", out)
	}
	if !strings.Contains(out, "password=***") {
		t.Errorf("expected redacted KV pair, got %q", out)
	}
	if !strings.Contains(out, "user=alice") {
		t.Errorf("non-sensitive KV should be preserved, got %q", out)
	}
}

func TestMaskKVNoMatch(t *testing.T) {
	m := New([]string{"token"}, "***")
	line := `level=info user=alice msg=ok`
	out := m.Mask(line)
	if out != line {
		t.Errorf("line without target field should be unchanged, got %q", out)
	}
	if m.MaskedCount() != 0 {
		t.Errorf("expected 0 masked lines, got %d", m.MaskedCount())
	}
}

func TestMaskedCountAccumulates(t *testing.T) {
	m := New([]string{"secret"}, "")
	lines := []string{
		`{"secret":"abc"}`,
		`level=info secret=xyz`,
		`level=info msg=ok`,
	}
	for _, l := range lines {
		m.Mask(l)
	}
	if m.MaskedCount() != 2 {
		t.Errorf("expected 2 masked lines, got %d", m.MaskedCount())
	}
}

func TestCustomPlaceholder(t *testing.T) {
	m := New([]string{"key"}, "<hidden>")
	line := `key=value other=ok`
	out := m.Mask(line)
	if !strings.Contains(out, "key=<hidden>") {
		t.Errorf("expected custom placeholder, got %q", out)
	}
}
