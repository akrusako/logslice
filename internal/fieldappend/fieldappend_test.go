package fieldappend

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestNewEmptyFieldReturnsError(t *testing.T) {
	_, err := New("", "_suffix")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewEmptySuffixReturnsError(t *testing.T) {
	_, err := New("level", "")
	if err == nil {
		t.Fatal("expected error for empty suffix")
	}
}

func TestNewValidParams(t *testing.T) {
	a, err := New("level", "_appended")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a == nil {
		t.Fatal("expected non-nil Appender")
	}
}

func TestAppendKVField(t *testing.T) {
	a, _ := New("env", "-prod")
	out := a.Apply("ts=2024-01-01 env=staging svc=api")
	if !strings.Contains(out, "env=staging-prod") {
		t.Errorf("expected appended value, got: %s", out)
	}
}

func TestKVFieldMissingNotModified(t *testing.T) {
	a, _ := New("missing", "_x")
	line := "ts=2024-01-01 level=info"
	out := a.Apply(line)
	if out != line {
		t.Errorf("expected unchanged line, got: %s", out)
	}
	if a.Missed() != 1 {
		t.Errorf("expected missed=1, got %d", a.Missed())
	}
}

func TestAppendJSONField(t *testing.T) {
	a, _ := New("service", ".v2")
	out := a.Apply(`{"service":"auth","level":"info"}`)
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(out), &m); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if m["service"] != "auth.v2" {
		t.Errorf("expected 'auth.v2', got %v", m["service"])
	}
}

func TestJSONFieldMissingNotModified(t *testing.T) {
	a, _ := New("missing", "_x")
	line := `{"level":"warn"}`
	out := a.Apply(line)
	var m map[string]interface{}
	json.Unmarshal([]byte(out), &m)
	if _, ok := m["missing"]; ok {
		t.Error("field should not have been injected")
	}
	if a.Missed() != 1 {
		t.Errorf("expected missed=1, got %d", a.Missed())
	}
}

func TestCountAccumulates(t *testing.T) {
	a, _ := New("k", "_s")
	for i := 0; i < 5; i++ {
		a.Apply("k=val other=x")
	}
	if a.Count() != 5 {
		t.Errorf("expected count=5, got %d", a.Count())
	}
}

func TestEmptyLineReturnsEmpty(t *testing.T) {
	a, _ := New("k", "_s")
	out := a.Apply("")
	if out != "" {
		t.Errorf("expected empty output, got %q", out)
	}
}
