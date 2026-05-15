package fieldtrim

import (
	"encoding/json"
	"testing"
)

func TestNewEmptyFieldReturnsError(t *testing.T) {
	_, err := New("", "")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewValidParams(t *testing.T) {
	tr, err := New("msg", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr == nil {
		t.Fatal("expected non-nil trimmer")
	}
}

func TestTrimKVWhitespace(t *testing.T) {
	tr, _ := New("msg", "")
	out := tr.Apply(`level=info msg="  hello  " host=web`)
	if want := `level=info msg=hello host=web`; out != want {
		t.Errorf("got %q, want %q", out, want)
	}
}

func TestTrimKVCustomCutset(t *testing.T) {
	tr, _ := New("path", "/")
	out := tr.Apply(`method=GET path=/api/v1/ status=200`)
	if want := `method=GET path=api/v1 status=200`; out != want {
		t.Errorf("got %q, want %q", out, want)
	}
}

func TestTrimJSONField(t *testing.T) {
	tr, _ := New("msg", "")
	out := tr.Apply(`{"level":"info","msg":"  hello  "}`)
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(out), &m); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if m["msg"] != "hello" {
		t.Errorf("got %q, want %q", m["msg"], "hello")
	}
}

func TestTrimJSONMissingFieldUnchanged(t *testing.T) {
	tr, _ := New("missing", "")
	line := `{"level":"info","msg":"  hello  "}`
	out := tr.Apply(line)
	if out != line {
		t.Errorf("expected unchanged, got %q", out)
	}
}

func TestEmptyLineUnchanged(t *testing.T) {
	tr, _ := New("msg", "")
	if out := tr.Apply(""); out != "" {
		t.Errorf("expected empty, got %q", out)
	}
}

func TestTrimmedCounterIncrements(t *testing.T) {
	tr, _ := New("msg", "")
	tr.Apply(`level=info msg="  hi  "`)
	tr.Apply(`level=info msg=clean`)
	if tr.Trimmed() != 1 {
		t.Errorf("expected 1 trimmed, got %d", tr.Trimmed())
	}
}

func TestKVFieldNotPresent(t *testing.T) {
	tr, _ := New("other", "")
	line := `level=info msg="  hello  "`
	out := tr.Apply(line)
	if out != line {
		t.Errorf("expected unchanged, got %q", out)
	}
}
