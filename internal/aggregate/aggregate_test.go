package aggregate_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/logslice/logslice/internal/aggregate"
)

func TestRecordNoField(t *testing.T) {
	c := aggregate.New("")
	c.Record(`level=info msg="hello"`)
	c.Record(`level=warn msg="bye"`)
	if c.Total() != 2 {
		t.Fatalf("expected total 2, got %d", c.Total())
	}
	counts := c.Counts()
	if counts["*"] != 2 {
		t.Errorf("expected wildcard bucket=2, got %d", counts["*"])
	}
}

func TestRecordKVField(t *testing.T) {
	c := aggregate.New("level")
	c.Record(`level=info msg="a"`)
	c.Record(`level=info msg="b"`)
	c.Record(`level=warn msg="c"`)
	counts := c.Counts()
	if counts["info"] != 2 {
		t.Errorf("expected info=2, got %d", counts["info"])
	}
	if counts["warn"] != 1 {
		t.Errorf("expected warn=1, got %d", counts["warn"])
	}
}

func TestRecordJSONField(t *testing.T) {
	c := aggregate.New("level")
	c.Record(`{"level":"error","msg":"boom"}`)
	c.Record(`{"level":"error","msg":"again"}`)
	counts := c.Counts()
	if counts["error"] != 2 {
		t.Errorf("expected error=2, got %d", counts["error"])
	}
}

func TestMissingFieldBucket(t *testing.T) {
	c := aggregate.New("level")
	c.Record("no fields here at all")
	counts := c.Counts()
	if counts["(missing)"] != 1 {
		t.Errorf("expected (missing)=1, got %d", counts["(missing)"])
	}
}

func TestWriteSummaryOrdered(t *testing.T) {
	c := aggregate.New("level")
	for i := 0; i < 3; i++ {
		c.Record(`level=error msg="e"`)
	}
	c.Record(`level=warn msg="w"`)

	var buf bytes.Buffer
	if err := c.WriteSummary(&buf); err != nil {
		t.Fatalf("WriteSummary error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	// highest count first
	if !strings.Contains(lines[0], "error") {
		t.Errorf("expected error first, got: %s", lines[0])
	}
}

func TestTotalAccumulates(t *testing.T) {
	c := aggregate.New("level")
	for i := 0; i < 10; i++ {
		c.Record(`level=debug`)
	}
	if c.Total() != 10 {
		t.Errorf("expected 10, got %d", c.Total())
	}
}
