package sample

import (
	"testing"
)

func TestRateOneKeepsAll(t *testing.T) {
	s := New(1.0, 42)
	for i := 0; i < 100; i++ {
		if !s.Keep("line") {
			t.Fatal("rate=1.0 should keep all lines")
		}
	}
	total, kept := s.Stats()
	if total != 100 || kept != 100 {
		t.Errorf("expected 100/100, got %d/%d", kept, total)
	}
}

func TestRateZeroDropsAll(t *testing.T) {
	s := New(0.0, 42)
	for i := 0; i < 100; i++ {
		if s.Keep("line") {
			t.Fatal("rate=0.0 should drop all lines")
		}
	}
	total, kept := s.Stats()
	if total != 100 || kept != 0 {
		t.Errorf("expected 0/100, got %d/%d", kept, total)
	}
}

func TestRateClampAboveOne(t *testing.T) {
	s := New(2.5, 1)
	if s.Rate() != 1.0 {
		t.Errorf("expected rate clamped to 1.0, got %f", s.Rate())
	}
}

func TestRateClampBelowZero(t *testing.T) {
	s := New(-0.5, 1)
	if s.Rate() != 0.0 {
		t.Errorf("expected rate clamped to 0.0, got %f", s.Rate())
	}
}

func TestPartialRateApproximate(t *testing.T) {
	s := New(0.5, 99)
	const n = 10000
	for i := 0; i < n; i++ {
		s.Keep("line")
	}
	_, kept := s.Stats()
	pct := float64(kept) / float64(n)
	if pct < 0.40 || pct > 0.60 {
		t.Errorf("expected ~50%% kept, got %.2f%%", pct*100)
	}
}

func TestStatsCountsTotal(t *testing.T) {
	s := New(0.5, 7)
	for i := 0; i < 20; i++ {
		s.Keep("x")
	}
	total, _ := s.Stats()
	if total != 20 {
		t.Errorf("expected total=20, got %d", total)
	}
}
