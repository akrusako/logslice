package fieldpresence

import (
	"testing"
)

func TestNewBothEmptyReturnsError(t *testing.T) {
	_, err := New(nil, nil)
	if err == nil {
		t.Fatal("expected error for empty require and forbid")
	}
}

func TestNewValidParams(t *testing.T) {
	c, err := New([]string{"level"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil Checker")
	}
}

func TestRequireFieldKV(t *testing.T) {
	c, _ := New([]string{"level"}, nil)
	if !c.Allow("level=info msg=hello") {
		t.Error("expected line with required field to pass")
	}
}

func TestRequireFieldMissingKV(t *testing.T) {
	c, _ := New([]string{"level"}, nil)
	if c.Allow("msg=hello") {
		t.Error("expected line missing required field to be dropped")
	}
}

func TestRequireFieldJSON(t *testing.T) {
	c, _ := New([]string{"level"}, nil)
	if !c.Allow(`{"level":"error","msg":"oops"}`) {
		t.Error("expected JSON line with required field to pass")
	}
}

func TestForbidFieldKV(t *testing.T) {
	c, _ := New(nil, []string{"debug"})
	if c.Allow("debug=true msg=verbose") {
		t.Error("expected line with forbidden field to be dropped")
	}
}

func TestForbidFieldAbsentKV(t *testing.T) {
	c, _ := New(nil, []string{"debug"})
	if !c.Allow("level=info msg=hello") {
		t.Error("expected line without forbidden field to pass")
	}
}

func TestCombinedRequireAndForbid(t *testing.T) {
	c, _ := New([]string{"level"}, []string{"trace"})
	if !c.Allow("level=info msg=ok") {
		t.Error("expected pass when required present and forbidden absent")
	}
	if c.Allow("level=info trace=true") {
		t.Error("expected drop when forbidden field present")
	}
	if c.Allow("msg=ok") {
		t.Error("expected drop when required field absent")
	}
}

func TestCounters(t *testing.T) {
	c, _ := New([]string{"level"}, nil)
	c.Allow("level=info")
	c.Allow("msg=only")
	if c.Checked() != 2 {
		t.Errorf("expected 2 checked, got %d", c.Checked())
	}
	if c.Passed() != 1 {
		t.Errorf("expected 1 passed, got %d", c.Passed())
	}
}
