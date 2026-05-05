package bsearch_test

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/bsearch"
	"github.com/yourorg/logslice/internal/lineindex"
	"github.com/yourorg/logslice/internal/timeparse"
)

// buildLog creates a multi-line log where each line starts with an RFC3339
// timestamp spaced one minute apart from baseTime.
func buildLog(base time.Time, count int) string {
	var sb strings.Builder
	for i := 0; i < count; i++ {
		t := base.Add(time.Duration(i) * time.Minute)
		fmt.Fprintf(&sb, "%s level=info msg=\"line %d\"\n", t.Format(time.RFC3339), i)
	}
	return sb.String()
}

func setup(t *testing.T, count int) (*bsearch.Finder, time.Time) {
	t.Helper()
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	log := buildLog(base, count)
	r := strings.NewReader(log)
	idx, err := lineindex.Build(r)
	if err != nil {
		t.Fatalf("lineindex.Build: %v", err)
	}
	_, _ = r.Seek(0, io.SeekStart)
	p := timeparse.NewParser()
	return bsearch.New(idx, bytes.NewReader([]byte(log)), p), base
}

func TestFindRangeFullSpan(t *testing.T) {
	f, base := setup(t, 10)
	first, last := f.FindRange(base, base.Add(9*time.Minute))
	if first != 0 || last != 9 {
		t.Errorf("expected [0,9], got [%d,%d]", first, last)
	}
}

func TestFindRangeSubset(t *testing.T) {
	f, base := setup(t, 10)
	start := base.Add(2 * time.Minute)
	end := base.Add(5 * time.Minute)
	first, last := f.FindRange(start, end)
	if first != 2 || last != 5 {
		t.Errorf("expected [2,5], got [%d,%d]", first, last)
	}
}

func TestFindRangeNoMatch(t *testing.T) {
	f, base := setup(t, 5)
	start := base.Add(10 * time.Minute)
	end := base.Add(20 * time.Minute)
	first, last := f.FindRange(start, end)
	if first != -1 || last != -1 {
		t.Errorf("expected [-1,-1], got [%d,%d]", first, last)
	}
}

func TestFindRangeEmptyIndex(t *testing.T) {
	p := timeparse.NewParser()
	idx, _ := lineindex.Build(strings.NewReader(""))
	f := bsearch.New(idx, bytes.NewReader(nil), p)
	base := time.Now()
	first, last := f.FindRange(base, base.Add(time.Hour))
	if first != -1 || last != -1 {
		t.Errorf("expected [-1,-1] for empty index, got [%d,%d]", first, last)
	}
}

func TestFindRangeSingleLine(t *testing.T) {
	f, base := setup(t, 1)
	first, last := f.FindRange(base, base)
	if first != 0 || last != 0 {
		t.Errorf("expected [0,0], got [%d,%d]", first, last)
	}
}
