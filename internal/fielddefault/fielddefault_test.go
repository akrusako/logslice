package fielddefault

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestNewEmptyFieldReturnsError(t *testing.T) {
	_, err := New("", "fallback")
	if err == nil {
		t.Fatal("expected error for empty field name")
	}
}

func TestNewValidParams(t *testing.T) {
	d, err := New("level", "info")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d == nil {
		t.Fatal("expected non-nil Defaulter")
	}
}

func TestApplyKVMissingField(t *testing.T) {
	d, _ := New("level", "info")
	out := d.Apply("msg=hello ts=2024-01-01")
	if !strings.Contains(out, "level=info") {
		t.Errorf("expected default injected, got: %s", out)
	}
	if d.Applied() != 1 {
		t.Errorf("expected applied=1, got %d", d.Applied())
	}
}

func TestApplyKVExistingFieldUnchanged(t *testing.T) {
	d, _ := New("level", "info")
	out := d.Apply("msg=hello level=warn")
	if !strings.Contains(out, "level=warn") {
		t.Errorf("expected original value preserved, got: %s", out)
	}
	if d.Skipped() != 1 {
		t.Errorf("expected skipped=1, got %d", d.Skipped())
	}
}

func TestApplyJSONMissingField(t *testing.T) {
	d, _ := New("level", "info")
	out := d.Apply(`{"msg":"hello"}`)
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(out), &m); err != nil {
		t.Fatalf("output not valid JSON: %v", err)
	}
	if m["level"] != "info" {
		t.Errorf("expected level=info in JSON, got %v", m["level"])
	}
}

func TestApplyJSONExistingFieldUnchanged(t *testing.T) {
	d, _ := New("level", "info")
	out := d.Apply(`{"level":"error","msg":"boom"}`)
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(out), &m); err != nil {
		t.Fatalf("output not valid JSON: %v", err)
	}
	if m["level"] != "error" {
		t.Errorf("expected original level preserved, got %v", m["level"])
	}
	if d.Skipped() != 1 {
		t.Errorf("expected skipped=1, got %d", d.Skipped())
	}
}

func TestEmptyLineReturnsEmpty(t *testing.T) {
	d, _ := New("level", "info")
	out := d.Apply("")
	if out != "" {
		t.Errorf("expected empty output, got %q", out)
	}
}

func TestCountersAccumulate(t *testing.T) {
	d, _ := New("env", "prod")
	d.Apply("msg=a")
	d.Apply("msg=b")
	d.Apply("msg=c env=staging")
	if d.Applied() != 2 {
		t.Errorf("expected applied=2, got %d", d.Applied())
	}
	if d.Skipped() != 1 {
		t.Errorf("expected skipped=1, got %d", d.Skipped())
	}
}
