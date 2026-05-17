package fieldformat

import (
	"strings"
	"testing"
)

func TestNewEmptyFieldReturnsError(t *testing.T) {
	_, err := New("", "%s")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewEmptyFormatReturnsError(t *testing.T) {
	_, err := New("level", "")
	if err == nil {
		t.Fatal("expected error for empty format")
	}
}

func TestNewValidParams(t *testing.T) {
	f, err := New("score", "%.2f")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil formatter")
	}
}

func TestFormatKVField(t *testing.T) {
	f, _ := New("val", "%05s")
	out := f.Apply("ts=2024 val=42 host=srv")
	if !strings.Contains(out, "val=   42") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatJSONField(t *testing.T) {
	f, _ := New("score", "%.1f")
	out := f.Apply(`{"score":3.14159,"host":"a"}`)
	if !strings.Contains(out, `"score":"3.1"`) {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestMissingFieldKVUnchanged(t *testing.T) {
	f, _ := New("missing", "%s")
	line := "ts=2024 host=srv"
	out := f.Apply(line)
	if out != line {
		t.Errorf("expected unchanged, got %q", out)
	}
}

func TestMissingFieldJSONUnchanged(t *testing.T) {
	f, _ := New("missing", "%s")
	line := `{"host":"srv"}`
	out := f.Apply(line)
	if out != line {
		t.Errorf("expected unchanged, got %q", out)
	}
}

func TestFormattedCountIncrements(t *testing.T) {
	f, _ := New("val", "%05s")
	f.Apply("val=1 host=a")
	f.Apply("val=2 host=b")
	f.Apply("host=c") // no field
	if f.FormattedCount() != 2 {
		t.Errorf("expected 2, got %d", f.FormattedCount())
	}
}

func TestEmptyLineReturnsEmpty(t *testing.T) {
	f, _ := New("val", "%s")
	out := f.Apply("")
	if out != "" {
		t.Errorf("expected empty, got %q", out)
	}
}
