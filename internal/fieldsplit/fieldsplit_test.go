package fieldsplit

import (
	"encoding/json"
	"testing"
)

func TestNewEmptyFieldReturnsError(t *testing.T) {
	_, err := New("", ",", "part_")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewEmptyDelimReturnsError(t *testing.T) {
	_, err := New("tags", "", "part_")
	if err == nil {
		t.Fatal("expected error for empty delim")
	}
}

func TestNewValidParams(t *testing.T) {
	s, err := New("tags", ",", "tag_")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil splitter")
	}
}

func TestDefaultPrefixDerivedFromField(t *testing.T) {
	s, err := New("tags", ",", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result := s.Apply("tags=a,b,c")
	if result == "tags=a,b,c" {
		t.Fatal("expected field to be split")
	}
}

func TestSplitKVField(t *testing.T) {
	s, _ := New("tags", ",", "tag_")
	result := s.Apply("level=info tags=a,b,c msg=hello")
	if !contains(result, "tag_0=a") || !contains(result, "tag_1=b") || !contains(result, "tag_2=c") {
		t.Errorf("unexpected result: %s", result)
	}
	if contains(result, "tags=") {
		t.Errorf("original field should be removed: %s", result)
	}
	if s.Split() != 1 {
		t.Errorf("expected split count 1, got %d", s.Split())
	}
}

func TestSplitKVMissingField(t *testing.T) {
	s, _ := New("tags", ",", "tag_")
	line := "level=info msg=hello"
	result := s.Apply(line)
	if result != line {
		t.Errorf("expected original line, got %s", result)
	}
	if s.Split() != 0 {
		t.Errorf("expected split count 0, got %d", s.Split())
	}
}

func TestSplitJSONField(t *testing.T) {
	s, _ := New("tags", ",", "tag_")
	line := `{"level":"info","tags":"a,b,c"}`
	result := s.Apply(line)
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(result), &m); err != nil {
		t.Fatalf("result is not valid JSON: %v", err)
	}
	if m["tag_0"] != "a" || m["tag_1"] != "b" || m["tag_2"] != "c" {
		t.Errorf("unexpected fields: %v", m)
	}
	if _, ok := m["tags"]; ok {
		t.Error("original field should be removed")
	}
	if s.Split() != 1 {
		t.Errorf("expected split count 1, got %d", s.Split())
	}
}

func TestSplitJSONMissingField(t *testing.T) {
	s, _ := New("tags", ",", "tag_")
	line := `{"level":"info"}`
	result := s.Apply(line)
	if result != line {
		t.Errorf("expected original line, got %s", result)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		len(s) > 0 && containsStr(s, substr))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
