package context_test

import (
	"context"
	"testing"
	"time"

	pctx "github.com/yourorg/logslice/internal/context"
)

func TestNewSetsStart(t *testing.T) {
	before := time.Now()
	pc := pctx.New(context.Background(), pctx.Config{})
	if pc.Start.Before(before) {
		t.Fatal("Start should be >= before")
	}
}

func TestRecordLineUnlimited(t *testing.T) {
	pc := pctx.New(context.Background(), pctx.Config{MaxLines: 0})
	for i := 0; i < 1000; i++ {
		if !pc.RecordLine() {
			t.Fatalf("unlimited pipeline stopped at line %d", i)
		}
	}
}

func TestRecordLineCapAtMax(t *testing.T) {
	pc := pctx.New(context.Background(), pctx.Config{MaxLines: 3})
	if !pc.RecordLine() { t.Fatal("line 1 should pass") }
	if !pc.RecordLine() { t.Fatal("line 2 should pass") }
	if pc.RecordLine() { t.Fatal("line 3 should signal stop") }
}

func TestLinesOut(t *testing.T) {
	pc := pctx.New(context.Background(), pctx.Config{MaxLines: 10})
	for i := 0; i < 4; i++ {
		pc.RecordLine()
	}
	if pc.LinesOut() != 4 {
		t.Fatalf("expected 4 lines out, got %d", pc.LinesOut())
	}
}

func TestInTimeRangeOpenBounds(t *testing.T) {
	pc := pctx.New(context.Background(), pctx.Config{})
	if !pc.InTimeRange(time.Now()) {
		t.Fatal("open bounds should accept any time")
	}
}

func TestInTimeRangeSince(t *testing.T) {
	now := time.Now()
	pc := pctx.New(context.Background(), pctx.Config{Since: now})
	if pc.InTimeRange(now.Add(-time.Second)) {
		t.Fatal("time before Since should be rejected")
	}
	if !pc.InTimeRange(now.Add(time.Second)) {
		t.Fatal("time after Since should be accepted")
	}
}

func TestInTimeRangeUntil(t *testing.T) {
	now := time.Now()
	pc := pctx.New(context.Background(), pctx.Config{Until: now})
	if pc.InTimeRange(now.Add(time.Second)) {
		t.Fatal("time after Until should be rejected")
	}
	if !pc.InTimeRange(now.Add(-time.Second)) {
		t.Fatal("time before Until should be accepted")
	}
}

func TestElapsedPositive(t *testing.T) {
	pc := pctx.New(context.Background(), pctx.Config{})
	time.Sleep(time.Millisecond)
	if pc.Elapsed() <= 0 {
		t.Fatal("elapsed should be positive after sleep")
	}
}
