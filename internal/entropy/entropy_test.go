package entropy

import (
	"math"
	"testing"
)

func TestScoreUniform(t *testing.T) {
	// A string of identical bytes has entropy 0.
	got := Score("aaaaaaaaaa")
	if got != 0 {
		t.Fatalf("expected 0, got %f", got)
	}
}

func TestScoreMaxEntropy(t *testing.T) {
	// Two equally-probable symbols → entropy == 1 bit/byte.
	got := Score("ababababab")
	if math.Abs(got-1.0) > 1e-9 {
		t.Fatalf("expected 1.0, got %f", got)
	}
}

func TestScoreEmptyString(t *testing.T) {
	if Score("") != 0 {
		t.Fatal("empty string should have entropy 0")
	}
}

func TestScoreNaturalText(t *testing.T) {
	// Natural log text should have moderate entropy (roughly 3–5 bits).
	s := `2024-01-15T12:00:00Z level=info msg="user logged in" user=alice`
	h := Score(s)
	if h < 3.0 || h > 6.0 {
		t.Fatalf("unexpected entropy for natural text: %f", h)
	}
}

func TestZeroThresholdAllowsAll(t *testing.T) {
	sc := New(0)
	lines := []string{"aaaaaa", "hello world", string(make([]byte, 256))}
	for _, l := range lines {
		if !sc.Allow(l) {
			t.Fatalf("zero threshold should allow all lines, blocked: %q", l)
		}
	}
	if sc.Above() != 0 {
		t.Fatalf("expected 0 above, got %d", sc.Above())
	}
}

func TestHighEntropyLineBlocked(t *testing.T) {
	sc := New(4.5)
	// Construct a high-entropy line (all 256 byte values).
	b := make([]byte, 256)
	for i := range b {
		b[i] = byte(i)
	}
	if sc.Allow(string(b)) {
		t.Fatal("high-entropy line should be blocked")
	}
	if sc.Above() != 1 {
		t.Fatalf("expected Above=1, got %d", sc.Above())
	}
}

func TestLowEntropyLineAllowed(t *testing.T) {
	sc := New(4.5)
	if !sc.Allow("aaabbbccc") {
		t.Fatal("low-entropy line should be allowed")
	}
	if sc.Above() != 0 {
		t.Fatalf("expected Above=0, got %d", sc.Above())
	}
}

func TestTotalCounter(t *testing.T) {
	sc := New(3.0)
	for i := 0; i < 5; i++ {
		sc.Allow("some log line")
	}
	if sc.Total() != 5 {
		t.Fatalf("expected Total=5, got %d", sc.Total())
	}
}
