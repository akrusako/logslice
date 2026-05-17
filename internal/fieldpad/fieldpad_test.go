package fieldpad

import (
	"strings"
	"testing"
)

func TestNewEmptyFieldReturnsError(t *testing.T) {
	_, err := New("", 8, ' ', Right)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewZeroWidthReturnsError(t *testing.T) {
	_, err := New("level", 0, ' ', Right)
	if err == nil {
		t.Fatal("expected error for zero width")
	}
}

func TestNewValidParams(t *testing.T) {
	p, err := New("level", 8, ' ', Right)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil Padder")
	}
}

func TestPadKVRightPad(t *testing.T) {
	p, _ := New("level", 8, ' ', Right)
	out := p.Apply("ts=2024-01-01 level=INFO msg=start")
	if !strings.Contains(out, "level=INFO    ") {
		t.Errorf("expected right-padded value, got: %s", out)
	}
}

func TestPadKVLeftPad(t *testing.T) {
	p, _ := New("level", 8, ' ', Left)
	out := p.Apply("ts=2024-01-01 level=INFO msg=start")
	if !strings.Contains(out, "level=    INFO") {
		t.Errorf("expected left-padded value, got: %s", out)
	}
}

func TestPadKVMissingFieldUnchanged(t *testing.T) {
	p, _ := New("severity", 8, ' ', Right)
	line := "ts=2024-01-01 level=INFO msg=start"
	out := p.Apply(line)
	if out != line {
		t.Errorf("expected unchanged line, got: %s", out)
	}
}

func TestPadJSONRightPad(t *testing.T) {
	p, _ := New("level", 8, '-', Right)
	out := p.Apply(`{"level":"INFO","msg":"ok"}`)
	if !strings.Contains(out, "INFO---") {
		t.Errorf("expected right-padded JSON value, got: %s", out)
	}
}

func TestPadJSONLeftPad(t *testing.T) {
	p, _ := New("level", 8, '0', Left)
	out := p.Apply(`{"level":"ERR","msg":"fail"}`)
	if !strings.Contains(out, "00000ERR") {
		t.Errorf("expected left-padded JSON value, got: %s", out)
	}
}

func TestAlreadyWideEnoughUnchanged(t *testing.T) {
	p, _ := New("level", 3, ' ', Right)
	line := "level=INFO msg=start"
	out := p.Apply(line)
	if out != line {
		t.Errorf("expected unchanged line when value already meets width, got: %s", out)
	}
}

func TestPaddedCountIncrements(t *testing.T) {
	p, _ := New("level", 8, ' ', Right)
	p.Apply("level=INFO msg=a")
	p.Apply("level=WARN msg=b")
	p.Apply("msg=no-level-field")
	if p.PaddedCount() != 2 {
		t.Errorf("expected 2 padded, got %d", p.PaddedCount())
	}
}

func TestDefaultPadCharIsSpace(t *testing.T) {
	p, _ := New("code", 6, 0, Right)
	out := p.Apply("code=OK msg=done")
	if !strings.Contains(out, "code=OK    ") {
		t.Errorf("expected space padding, got: %s", out)
	}
}
