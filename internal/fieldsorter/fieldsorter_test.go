package fieldsorter

import (
	"testing"
)

func TestNewEmptyOrderReturnsError(t *testing.T) {
	_, err := New([]string{})
	if err == nil {
		t.Fatal("expected error for empty order")
	}
}

func TestNewValidOrder(t *testing.T) {
	s, err := New([]string{"level", "msg"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil sorter")
	}
}

func TestSortKVLine(t *testing.T) {
	s, _ := New([]string{"level", "msg"})
	in := "msg=hello level=info ts=2024-01-01"
	out := s.Apply(in)
	if out != "level=info msg=hello ts=2024-01-01" {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestSortJSONLine(t *testing.T) {
	s, _ := New([]string{"level", "msg"})
	in := `{"ts":"2024-01-01","msg":"hello","level":"info"}`
	out := s.Apply(in)
	// Verify level comes before msg, both before ts
	levelIdx := indexOf(out, `"level"`)
	msgIdx := indexOf(out, `"msg"`)
	tsIdx := indexOf(out, `"ts"`)
	if levelIdx < 0 || msgIdx < 0 || tsIdx < 0 {
		t.Fatalf("missing fields in output: %q", out)
	}
	if levelIdx > msgIdx {
		t.Errorf("level should appear before msg, got: %q", out)
	}
	if msgIdx > tsIdx {
		t.Errorf("msg should appear before ts, got: %q", out)
	}
}

func TestSortedCounterIncrements(t *testing.T) {
	s, _ := New([]string{"level"})
	s.Apply("level=info msg=hello")
	s.Apply(`{"level":"warn","msg":"bye"}`)
	if s.Sorted() != 2 {
		t.Errorf("expected 2 sorted, got %d", s.Sorted())
	}
}

func TestSkippedCounterIncrements(t *testing.T) {
	s, _ := New([]string{"level"})
	s.Apply("plain text line")
	s.Apply("")
	if s.Skipped() != 2 {
		t.Errorf("expected 2 skipped, got %d", s.Skipped())
	}
}

func TestMissingFieldsIgnored(t *testing.T) {
	s, _ := New([]string{"level", "missing"})
	in := "level=info msg=hello"
	out := s.Apply(in)
	if out != "level=info msg=hello" {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestEmptyLineSkipped(t *testing.T) {
	s, _ := New([]string{"level"})
	out := s.Apply("")
	if out != "" {
		t.Fatalf("expected empty string, got %q", out)
	}
	if s.Skipped() != 1 {
		t.Errorf("expected 1 skipped")
	}
}

func indexOf(s, sub string) int {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
