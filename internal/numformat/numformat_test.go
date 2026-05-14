package numformat

import (
	"strings"
	"testing"
)

func TestNewEmptyFieldReturnsError(t *testing.T) {
	_, err := New("", StyleFixed, 2)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewNegativePrecClampedToZero(t *testing.T) {
	f, err := New("val", StyleFixed, -3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.prec != 0 {
		t.Fatalf("expected prec=0, got %d", f.prec)
	}
}

func TestFixedKV(t *testing.T) {
	f, _ := New("latency", StyleFixed, 3)
	out := f.Apply("ts=2024-01-01 latency=1.23456 status=ok")
	if !strings.Contains(out, "latency=1.235") {
		t.Fatalf("unexpected output: %s", out)
	}
}

func TestRoundKV(t *testing.T) {
	f, _ := New("score", StyleRound, 0)
	out := f.Apply("score=7.8 host=srv1")
	if !strings.Contains(out, "score=8") {
		t.Fatalf("unexpected output: %s", out)
	}
}

func TestSIKV(t *testing.T) {
	f, _ := New("bytes", StyleSI, 0)
	out := f.Apply("bytes=1500000 path=/log")
	if !strings.Contains(out, "bytes=1.50M") {
		t.Fatalf("unexpected output: %s", out)
	}
}

func TestFixedJSON(t *testing.T) {
	f, _ := New("duration", StyleFixed, 2)
	out := f.Apply(`{"duration":3.14159,"level":"info"}`)
	if !strings.Contains(out, `"duration":3.14`) {
		t.Fatalf("unexpected output: %s", out)
	}
}

func TestRoundJSON(t *testing.T) {
	f, _ := New("count", StyleRound, 0)
	out := f.Apply(`{"count":42.6,"svc":"api"}`)
	if !strings.Contains(out, `"count":43`) {
		t.Fatalf("unexpected output: %s", out)
	}
}

func TestSIJSON(t *testing.T) {
	f, _ := New("mem", StyleSI, 0)
	out := f.Apply(`{"mem":2000000000}`)
	if !strings.Contains(out, `"mem":"2.00G"`) {
		t.Fatalf("unexpected output: %s", out)
	}
}

func TestMissingFieldPassesThrough(t *testing.T) {
	f, _ := New("missing", StyleFixed, 2)
	line := "key=val other=123"
	out := f.Apply(line)
	if out != line {
		t.Fatalf("expected unchanged line, got: %s", out)
	}
	if f.Formatted() != 1 {
		t.Fatalf("expected formatted=1 (absent field is not an error), got %d", f.Formatted())
	}
}

func TestNonNumericFieldDropped(t *testing.T) {
	f, _ := New("val", StyleFixed, 2)
	out := f.Apply("val=notanumber extra=1")
	if out != "val=notanumber extra=1" {
		t.Fatalf("expected original line, got: %s", out)
	}
	if f.Dropped() != 1 {
		t.Fatalf("expected dropped=1, got %d", f.Dropped())
	}
}

func TestFormattedCounterIncrements(t *testing.T) {
	f, _ := New("v", StyleRound, 0)
	for i := 0; i < 5; i++ {
		f.Apply("v=1.5 x=y")
	}
	if f.Formatted() != 5 {
		t.Fatalf("expected 5, got %d", f.Formatted())
	}
}

func TestSISmallValue(t *testing.T) {
	f, _ := New("rate", StyleSI, 0)
	out := f.Apply("rate=0.5 unit=rps")
	if !strings.Contains(out, "rate=0.50") {
		t.Fatalf("unexpected output: %s", out)
	}
}
