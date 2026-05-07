package pipeline

import (
	"bytes"
	"strings"
	"testing"
)

func keepAll(line string) (string, bool)  { return line, true }
func dropAll(line string) (string, bool)  { return line, false }
func upper(line string) (string, bool)    { return strings.ToUpper(line), true }

func TestEmptyStagesPassThrough(t *testing.T) {
	var buf bytes.Buffer
	p := New(&buf)
	ok, err := p.Process("hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected line to pass")
	}
	if got := strings.TrimRight(buf.String(), "\n"); got != "hello" {
		t.Fatalf("want 'hello', got %q", got)
	}
}

func TestDropStageFiltersLine(t *testing.T) {
	var buf bytes.Buffer
	p := New(&buf, dropAll)
	ok, err := p.Process("hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected line to be dropped")
	}
	if buf.Len() != 0 {
		t.Fatalf("expected no output, got %q", buf.String())
	}
}

func TestTransformStageApplied(t *testing.T) {
	var buf bytes.Buffer
	p := New(&buf, upper)
	p.Process("hello")
	if got := strings.TrimRight(buf.String(), "\n"); got != "HELLO" {
		t.Fatalf("want 'HELLO', got %q", got)
	}
}

func TestStagesAppliedInOrder(t *testing.T) {
	var buf bytes.Buffer
	// upper then drop: line should be transformed then dropped
	p := New(&buf, upper, dropAll)
	ok, _ := p.Process("hello")
	if ok {
		t.Fatal("expected line to be dropped after upper")
	}
}

func TestStatsAccumulate(t *testing.T) {
	var buf bytes.Buffer
	p := New(&buf, keepAll)
	for i := 0; i < 5; i++ {
		p.Process("line")
	}
	p2 := New(&buf, dropAll)
	for i := 0; i < 3; i++ {
		p2.Process("line")
	}
	passed, dropped := p.Stats()
	if passed != 5 {
		t.Fatalf("want passed=5, got %d", passed)
	}
	if dropped != 0 {
		t.Fatalf("want dropped=0, got %d", dropped)
	}
	_, dropped2 := p2.Stats()
	if dropped2 != 3 {
		t.Fatalf("want dropped=3, got %d", dropped2)
	}
}
