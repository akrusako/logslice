package fieldclip

import (
	"testing"
)

func TestNewEmptyFieldReturnsError(t *testing.T) {
	_, err := New("", 0, 10)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewMinGreaterThanMaxReturnsError(t *testing.T) {
	_, err := New("val", 10, 5)
	if err == nil {
		t.Fatal("expected error when min > max")
	}
}

func TestNewValidParams(t *testing.T) {
	c, err := New("val", 0, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil clipper")
	}
}

func TestClampKVBelow(t *testing.T) {
	c, _ := New("score", 0, 100)
	out := c.Apply("score=-5 msg=hello")
	if out != "score=0 msg=hello" {
		t.Errorf("unexpected output: %q", out)
	}
	if c.Clipped() != 1 {
		t.Errorf("expected 1 clipped, got %d", c.Clipped())
	}
}

func TestClampKVAbove(t *testing.T) {
	c, _ := New("score", 0, 100)
	out := c.Apply("score=200 msg=hello")
	if out != "score=100 msg=hello" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestClampKVInRange(t *testing.T) {
	c, _ := New("score", 0, 100)
	line := "score=50 msg=hello"
	out := c.Apply(line)
	if out != line {
		t.Errorf("expected unchanged, got %q", out)
	}
	if c.Clipped() != 0 {
		t.Errorf("expected 0 clipped, got %d", c.Clipped())
	}
}

func TestClampJSONBelow(t *testing.T) {
	c, _ := New("temp", -10, 50)
	out := c.Apply(`{"temp":-99,"host":"a"}`)
	// value should be clamped to -10
	if out == `{"temp":-99,"host":"a"}` {
		t.Error("expected value to be clamped")
	}
	if c.Clipped() != 1 {
		t.Errorf("expected 1 clipped, got %d", c.Clipped())
	}
}

func TestClampJSONAbove(t *testing.T) {
	c, _ := New("temp", -10, 50)
	out := c.Apply(`{"temp":999,"host":"b"}`)
	if out == `{"temp":999,"host":"b"}` {
		t.Error("expected value to be clamped")
	}
}

func TestClampJSONMissingFieldUnchanged(t *testing.T) {
	c, _ := New("temp", 0, 100)
	line := `{"host":"x"}`
	out := c.Apply(line)
	if out != line {
		t.Errorf("expected unchanged, got %q", out)
	}
}

func TestPlainTextUnchanged(t *testing.T) {
	c, _ := New("val", 0, 10)
	line := "just a plain log line"
	out := c.Apply(line)
	if out != line {
		t.Errorf("expected unchanged, got %q", out)
	}
}

func TestClippedCountAccumulates(t *testing.T) {
	c, _ := New("n", 1, 5)
	c.Apply("n=0")
	c.Apply("n=6")
	c.Apply("n=3")
	if c.Clipped() != 2 {
		t.Errorf("expected 2 clipped, got %d", c.Clipped())
	}
}
