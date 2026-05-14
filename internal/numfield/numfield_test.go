package numfield_test

import (
	"testing"

	"github.com/user/logslice/internal/numfield"
)

func TestNewEmptyFieldReturnsError(t *testing.T) {
	_, err := numfield.New("", numfield.OpGT, 0)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestParseOpValid(t *testing.T) {
	cases := []string{">", ">=", "<", "<=", "=", "=="}
	for _, c := range cases {
		if _, err := numfield.ParseOp(c); err != nil {
			t.Errorf("ParseOp(%q) unexpected error: %v", c, err)
		}
	}
}

func TestParseOpInvalid(t *testing.T) {
	if _, err := numfield.ParseOp("!="); err == nil {
		t.Fatal("expected error for unknown operator")
	}
}

func TestAllowKVGreaterThan(t *testing.T) {
	f, _ := numfield.New("latency", numfield.OpGT, 100)
	if !f.Allow(`level=info latency=200 msg=slow`) {
		t.Error("expected line with latency=200 to pass GT 100")
	}
	if f.Allow(`level=info latency=50 msg=fast`) {
		t.Error("expected line with latency=50 to be dropped by GT 100")
	}
}

func TestAllowJSONLessThanOrEqual(t *testing.T) {
	f, _ := numfield.New("score", numfield.OpLTE, 5.0)
	if !f.Allow(`{"score":5,"msg":"ok"}`) {
		t.Error("expected score=5 to pass LTE 5")
	}
	if f.Allow(`{"score":6,"msg":"high"}`) {
		t.Error("expected score=6 to be dropped by LTE 5")
	}
}

func TestAllowEqualOperator(t *testing.T) {
	f, _ := numfield.New("code", numfield.OpEQ, 200)
	if !f.Allow(`{"code":200,"msg":"ok"}`) {
		t.Error("expected code=200 to pass EQ 200")
	}
	if f.Allow(`{"code":404,"msg":"not found"}`) {
		t.Error("expected code=404 to be dropped by EQ 200")
	}
}

func TestMissingFieldPassesThrough(t *testing.T) {
	f, _ := numfield.New("latency", numfield.OpGT, 100)
	if !f.Allow(`level=info msg=no-latency-here`) {
		t.Error("expected line without field to pass through")
	}
}

func TestNonNumericFieldPassesThrough(t *testing.T) {
	f, _ := numfield.New("latency", numfield.OpGT, 100)
	if !f.Allow(`latency=fast msg=non-numeric`) {
		t.Error("expected non-numeric field value to pass through")
	}
}

func TestCounters(t *testing.T) {
	f, _ := numfield.New("n", numfield.OpGTE, 3)
	lines := []string{
		`n=5 msg=a`, // pass
		`n=1 msg=b`, // drop
		`n=3 msg=c`, // pass
		`msg=d`,     // pass (missing field)
	}
	for _, l := range lines {
		f.Allow(l)
	}
	if f.Matched() != 3 {
		t.Errorf("matched: got %d, want 3", f.Matched())
	}
	if f.Dropped() != 1 {
		t.Errorf("dropped: got %d, want 1", f.Dropped())
	}
}
