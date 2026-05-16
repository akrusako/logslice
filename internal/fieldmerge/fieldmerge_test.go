package fieldmerge

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestNewEmptySrcAReturnsError(t *testing.T) {
	_, err := New("", "b", "dst", "-")
	if err == nil {
		t.Fatal("expected error for empty srcA")
	}
}

func TestNewEmptySrcBReturnsError(t *testing.T) {
	_, err := New("a", "", "dst", "-")
	if err == nil {
		t.Fatal("expected error for empty srcB")
	}
}

func TestNewEmptyDstReturnsError(t *testing.T) {
	_, err := New("a", "b", "", "-")
	if err == nil {
		t.Fatal("expected error for empty dst")
	}
}

func TestNewDstSameAsSrcReturnsError(t *testing.T) {
	_, err := New("a", "b", "a", "-")
	if err == nil {
		t.Fatal("expected error when dst equals srcA")
	}
}

func TestNewValidParams(t *testing.T) {
	m, err := New("first", "last", "full", " ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m == nil {
		t.Fatal("expected non-nil Merger")
	}
}

func TestMergeKVLine(t *testing.T) {
	m, _ := New("first", "last", "full", " ")
	out := m.Apply("first=John last=Doe level=info")
	if !strings.Contains(out, "full=John Doe") {
		t.Errorf("expected merged field in output, got: %s", out)
	}
	if m.Merged() != 1 {
		t.Errorf("expected Merged()=1, got %d", m.Merged())
	}
}

func TestMergeKVMissingFieldSkipped(t *testing.T) {
	m, _ := New("first", "last", "full", " ")
	out := m.Apply("first=John level=info")
	if out != "first=John level=info" {
		t.Errorf("expected unchanged line, got: %s", out)
	}
	if m.Skipped() != 1 {
		t.Errorf("expected Skipped()=1, got %d", m.Skipped())
	}
}

func TestMergeJSONLine(t *testing.T) {
	m, _ := New("host", "port", "addr", ":")
	out := m.Apply(`{"host":"localhost","port":"8080","level":"info"}`)
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(out), &obj); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if obj["addr"] != "localhost:8080" {
		t.Errorf("expected addr=localhost:8080, got %v", obj["addr"])
	}
	if m.Merged() != 1 {
		t.Errorf("expected Merged()=1, got %d", m.Merged())
	}
}

func TestMergeJSONMissingFieldSkipped(t *testing.T) {
	m, _ := New("host", "port", "addr", ":")
	line := `{"host":"localhost","level":"info"}`
	out := m.Apply(line)
	if out != line {
		t.Errorf("expected unchanged line, got: %s", out)
	}
	if m.Skipped() != 1 {
		t.Errorf("expected Skipped()=1, got %d", m.Skipped())
	}
}

func TestMergeEmptyLineReturnsEmpty(t *testing.T) {
	m, _ := New("a", "b", "c", "-")
	out := m.Apply("")
	if out != "" {
		t.Errorf("expected empty output, got: %q", out)
	}
}
