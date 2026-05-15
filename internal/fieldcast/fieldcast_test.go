package fieldcast

import (
	"testing"
)

func TestParseTypeValid(t *testing.T) {
	for _, tc := range []struct{ in, want string }{
		{"string", "string"},
		{"int", "int"},
		{"float", "float"},
		{"bool", "bool"},
		{"INT", "int"},
	} {
		got, err := ParseType(tc.in)
		if err != nil {
			t.Fatalf("ParseType(%q) unexpected error: %v", tc.in, err)
		}
		if string(got) != tc.want {
			t.Errorf("ParseType(%q) = %q; want %q", tc.in, got, tc.want)
		}
	}
}

func TestParseTypeInvalid(t *testing.T) {
	_, err := ParseType("bytes")
	if err == nil {
		t.Fatal("expected error for unknown type")
	}
}

func TestNewEmptyFieldReturnsError(t *testing.T) {
	_, err := New("", TypeInt)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestCastKVIntFromFloat(t *testing.T) {
	c, _ := New("latency", TypeInt)
	got := c.Apply("level=info latency=3.7 msg=ok")
	if got != "level=info latency=3 msg=ok" {
		t.Errorf("unexpected output: %q", got)
	}
	if c.Casted() != 1 {
		t.Errorf("casted counter = %d; want 1", c.Casted())
	}
}

func TestCastKVFloat(t *testing.T) {
	c, _ := New("score", TypeFloat)
	got := c.Apply("score=42 msg=test")
	if got != "score=42 msg=test" {
		t.Errorf("unexpected output: %q", got)
	}
}

func TestCastKVBool(t *testing.T) {
	c, _ := New("ok", TypeBool)
	got := c.Apply("ok=1 msg=done")
	if got != "ok=true msg=done" {
		t.Errorf("unexpected output: %q", got)
	}
}

func TestCastMissingFieldUnchanged(t *testing.T) {
	c, _ := New("missing", TypeInt)
	line := "level=info msg=hello"
	got := c.Apply(line)
	if got != line {
		t.Errorf("expected unchanged line, got %q", got)
	}
	if c.Casted() != 0 {
		t.Errorf("casted counter should be 0")
	}
}

func TestCastFailureCountsAndReturnsOriginal(t *testing.T) {
	c, _ := New("flag", TypeInt)
	line := "flag=notanumber msg=x"
	got := c.Apply(line)
	if got != line {
		t.Errorf("expected unchanged line on failure, got %q", got)
	}
	if c.Failed() != 1 {
		t.Errorf("failed counter = %d; want 1", c.Failed())
	}
}

func TestCastJSONField(t *testing.T) {
	c, _ := New("code", TypeInt)
	line := `{"code":"200","msg":"ok"}`
	got := c.Apply(line)
	if got != `{"code":200,"msg":"ok"}` {
		t.Errorf("unexpected JSON output: %q", got)
	}
}
