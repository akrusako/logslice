package logscanner_test

import (
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/logscanner"
)

const sampleLog = `2024-01-15T10:00:00Z INFO  service started
2024-01-15T10:01:00Z DEBUG request received id=42
2024-01-15T10:02:00Z ERROR disk full path=/var/log
2024-01-15T10:03:00Z INFO  recovered
no-timestamp-line
2024-01-15T10:04:00Z INFO  shutdown
`

func TestScanAll(t *testing.T) {
	s := logscanner.New(strings.NewReader(sampleLog), logscanner.Options{})
	var count int
	for s.Scan() {
		count++
	}
	if err := s.Err(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 6 {
		t.Fatalf("expected 6 lines, got %d", count)
	}
}

func TestScanLineNumbers(t *testing.T) {
	s := logscanner.New(strings.NewReader(sampleLog), logscanner.Options{})
	var last int64
	for s.Scan() {
		last = s.Line().Number
	}
	if last != 6 {
		t.Fatalf("expected last line number 6, got %d", last)
	}
}

func TestScanTimeRange(t *testing.T) {
	after := time.Date(2024, 1, 15, 10, 1, 0, 0, time.UTC)
	before := time.Date(2024, 1, 15, 10, 2, 59, 0, time.UTC)

	s := logscanner.New(strings.NewReader(sampleLog), logscanner.Options{
		After:  after,
		Before: before,
	})

	var lines []logscanner.Line
	for s.Scan() {
		lines = append(lines, s.Line())
	}
	if err := s.Err(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Expect lines at 10:01, 10:02 plus the no-timestamp line (not filtered)
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines in range, got %d", len(lines))
	}
}

func TestNoTimestampLinePassthrough(t *testing.T) {
	input := "no timestamp here\n"
	s := logscanner.New(strings.NewReader(input), logscanner.Options{
		After: time.Now().Add(time.Hour), // future — would filter timed lines
	})
	var count int
	for s.Scan() {
		count++
		if s.Line().HasTime {
			t.Fatal("expected HasTime=false")
		}
	}
	if count != 1 {
		t.Fatalf("expected 1 line, got %d", count)
	}
}
