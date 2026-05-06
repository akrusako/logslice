package stats

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestNewSetsStartTime(t *testing.T) {
	before := time.Now()
	s := New()
	after := time.Now()
	if s.StartTime.Before(before) || s.StartTime.After(after) {
		t.Errorf("StartTime %v not in expected range [%v, %v]", s.StartTime, before, after)
	}
}

func TestRecordLineMatched(t *testing.T) {
	s := New()
	s.RecordLine(true, 42)
	if s.LinesScanned != 1 {
		t.Errorf("expected LinesScanned=1, got %d", s.LinesScanned)
	}
	if s.LinesMatched != 1 {
		t.Errorf("expected LinesMatched=1, got %d", s.LinesMatched)
	}
	if s.LinesDropped != 0 {
		t.Errorf("expected LinesDropped=0, got %d", s.LinesDropped)
	}
	if s.BytesRead != 42 {
		t.Errorf("expected BytesRead=42, got %d", s.BytesRead)
	}
}

func TestRecordLineDropped(t *testing.T) {
	s := New()
	s.RecordLine(false, 10)
	if s.LinesDropped != 1 {
		t.Errorf("expected LinesDropped=1, got %d", s.LinesDropped)
	}
	if s.LinesMatched != 0 {
		t.Errorf("expected LinesMatched=0, got %d", s.LinesMatched)
	}
}

func TestMatchRateZeroWhenNoLines(t *testing.T) {
	s := New()
	if s.MatchRate() != 0 {
		t.Errorf("expected MatchRate=0 with no lines, got %f", s.MatchRate())
	}
}

func TestMatchRateCalculation(t *testing.T) {
	s := New()
	s.RecordLine(true, 1)
	s.RecordLine(true, 1)
	s.RecordLine(false, 1)
	s.RecordLine(false, 1)
	if got := s.MatchRate(); got != 0.5 {
		t.Errorf("expected MatchRate=0.5, got %f", got)
	}
}

func TestFinishAndElapsed(t *testing.T) {
	s := New()
	time.Sleep(2 * time.Millisecond)
	s.Finish()
	if s.Elapsed() < time.Millisecond {
		t.Errorf("expected elapsed >= 1ms, got %v", s.Elapsed())
	}
}

func TestPrintOutput(t *testing.T) {
	s := New()
	s.RecordLine(true, 100)
	s.RecordLine(false, 50)
	s.Finish()
	var buf bytes.Buffer
	s.Print(&buf)
	out := buf.String()
	for _, want := range []string{"scanned=2", "matched=1", "dropped=1", "bytes=150", "match_rate="} {
		if !strings.Contains(out, want) {
			t.Errorf("Print output missing %q; got: %s", want, out)
		}
	}
}
