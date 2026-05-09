package prefix

import "testing"

func TestEmptyPrefixReturnsOriginal(t *testing.T) {
	p := New("")
	got := p.Apply("hello world")
	if got != "hello world" {
		t.Fatalf("expected %q, got %q", "hello world", got)
	}
}

func TestPrefixIsApplied(t *testing.T) {
	p := New("[INFO] ")
	got := p.Apply("server started")
	want := "[INFO] server started"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestEmptyLineIsSkipped(t *testing.T) {
	p := New("[INFO] ")
	got := p.Apply("")
	if got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
	if p.Dropped() != 1 {
		t.Fatalf("expected dropped=1, got %d", p.Dropped())
	}
	if p.Count() != 0 {
		t.Fatalf("expected count=0, got %d", p.Count())
	}
}

func TestCountAccumulates(t *testing.T) {
	p := New(">> ")
	for i := 0; i < 5; i++ {
		p.Apply("line")
	}
	if p.Count() != 5 {
		t.Fatalf("expected count=5, got %d", p.Count())
	}
}

func TestSetPrefixChangesPrefix(t *testing.T) {
	p := New("A: ")
	p.SetPrefix("B: ")
	got := p.Apply("test")
	if got != "B: test" {
		t.Fatalf("expected %q, got %q", "B: test", got)
	}
}

func TestStripFoundPrefix(t *testing.T) {
	p := New("[DEBUG] ")
	stripped, ok := p.Strip("[DEBUG] something happened")
	if !ok {
		t.Fatal("expected prefix to be found")
	}
	if stripped != "something happened" {
		t.Fatalf("expected %q, got %q", "something happened", stripped)
	}
}

func TestStripMissingPrefix(t *testing.T) {
	p := New("[DEBUG] ")
	stripped, ok := p.Strip("no prefix here")
	if ok {
		t.Fatal("expected prefix not to be found")
	}
	if stripped != "no prefix here" {
		t.Fatalf("expected original line back, got %q", stripped)
	}
}

func TestStripEmptyPrefix(t *testing.T) {
	p := New("")
	stripped, ok := p.Strip("anything")
	if !ok {
		t.Fatal("empty prefix should always match")
	}
	if stripped != "anything" {
		t.Fatalf("expected %q, got %q", "anything", stripped)
	}
}
