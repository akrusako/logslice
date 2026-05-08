package timestamp

import (
	"testing"
	"time"
)

func TestExtractJSONTimestamp(t *testing.T) {
	e := New(nil, time.RFC3339, false)
	line := `{"time":"2024-01-15T10:30:00Z","msg":"hello"}`
	ts, out, ok := e.Extract(line)
	if !ok {
		t.Fatal("expected timestamp to be found")
	}
	if ts != "2024-01-15T10:30:00Z" {
		t.Fatalf("unexpected ts %q", ts)
	}
	if out != line {
		t.Fatalf("line should be unchanged without strip, got %q", out)
	}
}

func TestExtractKVTimestamp(t *testing.T) {
	e := New(nil, time.RFC3339, false)
	line := `time=2024-06-01T00:00:00Z level=info msg=boot`
	ts, _, ok := e.Extract(line)
	if !ok {
		t.Fatal("expected timestamp to be found")
	}
	if ts != "2024-06-01T00:00:00Z" {
		t.Fatalf("unexpected ts %q", ts)
	}
}

func TestExtractMissingField(t *testing.T) {
	e := New(nil, time.RFC3339, false)
	_, out, ok := e.Extract("no timestamp here")
	if ok {
		t.Fatal("expected no timestamp")
	}
	if out != "no timestamp here" {
		t.Fatalf("line should be returned unchanged, got %q", out)
	}
}

func TestExtractCustomField(t *testing.T) {
	e := New([]string{"logged_at"}, time.RFC3339, false)
	line := `logged_at=2024-03-10T08:00:00Z svc=api`
	_, _, ok := e.Extract(line)
	if !ok {
		t.Fatal("expected timestamp in custom field")
	}
}

func TestExtractStripsField(t *testing.T) {
	e := New(nil, time.RFC3339, true)
	line := `time=2024-01-01T00:00:00Z level=warn msg=test`
	_, out, ok := e.Extract(line)
	if !ok {
		t.Fatal("expected timestamp found")
	}
	if out == line {
		t.Fatal("expected field to be stripped from line")
	}
	if contains(out, "time=2024-01-01T00:00:00Z") {
		t.Fatalf("timestamp field still present in stripped output: %q", out)
	}
}

func TestCustomOutputFormat(t *testing.T) {
	e := New(nil, "2006-01-02", false)
	line := `{"time":"2024-07-04T12:00:00Z","msg":"fireworks"}`
	ts, _, ok := e.Extract(line)
	if !ok {
		t.Fatal("expected timestamp found")
	}
	if ts != "2024-07-04" {
		t.Fatalf("unexpected formatted ts %q", ts)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStr(s, sub))
}

func containsStr(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
