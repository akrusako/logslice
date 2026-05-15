package fieldrename_test

import (
	"encoding/json"
	"testing"

	"github.com/logslice/logslice/internal/fieldrename"
)

func TestNewInvalidPairReturnsError(t *testing.T) {
	_, err := fieldrename.New([]string{"noequalssign"})
	if err == nil {
		t.Fatal("expected error for missing '=' in pair")
	}
}

func TestNewEmptyKeyReturnsError(t *testing.T) {
	_, err := fieldrename.New([]string{"=newname"})
	if err == nil {
		t.Fatal("expected error for empty old key")
	}
}

func TestNewEmptyValueReturnsError(t *testing.T) {
	_, err := fieldrename.New([]string{"oldname="})
	if err == nil {
		t.Fatal("expected error for empty new key")
	}
}

func TestNoMappingReturnsOriginal(t *testing.T) {
	r, _ := fieldrename.New(nil)
	line := "level=info msg=hello"
	if got := r.Apply(line); got != line {
		t.Fatalf("expected %q, got %q", line, got)
	}
}

func TestRenameKVField(t *testing.T) {
	r, err := fieldrename.New([]string{"msg=message"})
	if err != nil {
		t.Fatal(err)
	}
	got := r.Apply("level=info msg=hello ts=2024-01-01")
	if !containsKV(got, "message", "hello") {
		t.Fatalf("expected 'message=hello' in %q", got)
	}
	if containsKey(got, "msg=") {
		t.Fatalf("old key 'msg' should be gone in %q", got)
	}
}

func TestRenameJSONField(t *testing.T) {
	r, err := fieldrename.New([]string{"msg=message"})
	if err != nil {
		t.Fatal(err)
	}
	got := r.Apply(`{"level":"info","msg":"hello"}`)
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(got), &m); err != nil {
		t.Fatalf("result not valid JSON: %v", err)
	}
	if _, ok := m["message"]; !ok {
		t.Fatalf("expected 'message' key in %v", m)
	}
	if _, ok := m["msg"]; ok {
		t.Fatalf("old key 'msg' should be absent in %v", m)
	}
}

func TestRenameKVNoMatchUnchanged(t *testing.T) {
	r, _ := fieldrename.New([]string{"missing=other"})
	line := "level=info msg=hello"
	if got := r.Apply(line); got != line {
		t.Fatalf("expected unchanged line, got %q", got)
	}
}

func TestRenamedCounterIncrements(t *testing.T) {
	r, _ := fieldrename.New([]string{"msg=message"})
	r.Apply("level=info msg=hello")
	r.Apply("level=warn msg=bye")
	r.Apply("level=debug note=nothing")
	if r.Renamed() != 2 {
		t.Fatalf("expected 2 renamed, got %d", r.Renamed())
	}
	if r.Processed() != 3 {
		t.Fatalf("expected 3 processed, got %d", r.Processed())
	}
}

func containsKV(line, key, value string) bool {
	return len(line) > 0 && (len(key+"="+value) > 0) &&
		(len(line) >= len(key)+1) &&
		(func() bool {
			for _, part := range splitFields(line) {
				if part == key+"="+value {
					return true
				}
			}
			return false
		})()
}

func containsKey(line, prefix string) bool {
	for _, part := range splitFields(line) {
		if len(part) >= len(prefix) && part[:len(prefix)] == prefix {
			return true
		}
	}
	return false
}

func splitFields(line string) []string {
	var out []string
	for _, f := range splitSpaces(line) {
		if f != "" {
			out = append(out, f)
		}
	}
	return out
}

func splitSpaces(s string) []string {
	var parts []string
	start := -1
	for i, c := range s {
		if c == ' ' || c == '\t' {
			if start >= 0 {
				parts = append(parts, s[start:i])
				start = -1
			}
		} else if start < 0 {
			start = i
		}
	}
	if start >= 0 {
		parts = append(parts, s[start:])
	}
	return parts
}
