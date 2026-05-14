package labelinject

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestNewEmptyPairsOK(t *testing.T) {
	inj, err := New(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inj == nil {
		t.Fatal("expected non-nil injector")
	}
}

func TestNewInvalidPairReturnsError(t *testing.T) {
	_, err := New([]string{"noequals"})
	if err == nil {
		t.Fatal("expected error for invalid pair")
	}
}

func TestNewEmptyKeyReturnsError(t *testing.T) {
	_, err := New([]string{"=value"})
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestNoLabelsReturnsOriginal(t *testing.T) {
	inj, _ := New(nil)
	line := "some plain log line"
	if got := inj.Apply(line); got != line {
		t.Fatalf("expected %q, got %q", line, got)
	}
}

func TestInjectIntoKVLine(t *testing.T) {
	inj, _ := New([]string{"env=prod", "region=us-east"})
	result := inj.Apply("level=info msg=hello")
	if !strings.Contains(result, "env=prod") {
		t.Errorf("missing env label in %q", result)
	}
	if !strings.Contains(result, "region=us-east") {
		t.Errorf("missing region label in %q", result)
	}
}

func TestInjectIntoJSONLine(t *testing.T) {
	inj, _ := New([]string{"env=staging"})
	result := inj.Apply(`{"level":"warn","msg":"oops"}`)
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(result), &m); err != nil {
		t.Fatalf("result is not valid JSON: %v", err)
	}
	if m["env"] != "staging" {
		t.Errorf("expected env=staging in JSON, got %v", m["env"])
	}
	if m["level"] != "warn" {
		t.Errorf("original field lost, got %v", m["level"])
	}
}

func TestInjectIntoPlainTextLine(t *testing.T) {
	inj, _ := New([]string{"host=box1"})
	result := inj.Apply("plain text log entry")
	if !strings.HasSuffix(result, "host=box1") {
		t.Errorf("expected suffix host=box1 in %q", result)
	}
}

func TestInjectedCounter(t *testing.T) {
	inj, _ := New([]string{"k=v"})
	for i := 0; i < 5; i++ {
		inj.Apply("msg=hello")
	}
	if inj.Injected() != 5 {
		t.Errorf("expected 5 injected, got %d", inj.Injected())
	}
}

func TestDuplicateKeyLastWins(t *testing.T) {
	inj, _ := New([]string{"env=dev", "env=prod"})
	result := inj.Apply("msg=test")
	if !strings.Contains(result, "env=prod") {
		t.Errorf("expected last value to win in %q", result)
	}
}
