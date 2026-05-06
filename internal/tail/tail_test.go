package tail_test

import (
	"os"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/tail"
)

func writeTempLog(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "tail-*.log")
	if err != nil {
		t.Fatalf("create temp: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestTailLessThanTotal(t *testing.T) {
	path := writeTempLog(t, "line1\nline2\nline3\nline4\nline5\n")
	r, err := tail.New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	lines, err := r.Lines(3)
	if err != nil {
		t.Fatalf("Lines: %v", err)
	}
	if len(lines) != 3 {
		t.Fatalf("want 3 lines, got %d", len(lines)))
	}
	if lines[0] != "line3" || lines[2] != "line5" {
		t.Errorf("unexpected lines: %v", lines)
	}
}

func TestTailMoreThanTotal(t *testing.T) {
	path := writeTempLog(t, "alpha\nbeta\n")
	r, err := tail.New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	lines, err := r.Lines(100)
	if err != nil {
		t.Fatalf("Lines: %v", err)
	}
	if len(lines) != 2 {
		t.Fatalf("want 2 lines, got %d", len(lines))
	}
}

func TestTailZero(t *testing.T) {
	path := writeTempLog(t, "only\n")
	r, err := tail.New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	lines, err := r.Lines(0)
	if err != nil {
		t.Fatalf("Lines: %v", err)
	}
	if len(lines) != 0 {
		t.Errorf("expected empty, got %v", lines)
	}
}

func TestTailContentIntegrity(t *testing.T) {
	input := strings.Repeat("2024-01-01T00:00:00Z level=info msg=hello\n", 20)
	path := writeTempLog(t, input)
	r, err := tail.New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	lines, err := r.Lines(5)
	if err != nil {
		t.Fatalf("Lines: %v", err)
	}
	for _, l := range lines {
		if !strings.Contains(l, "msg=hello") {
			t.Errorf("unexpected line content: %q", l)
		}
	}
}

// TestTailEmptyFile verifies that tailing an empty file returns no lines
// rather than an error, since an empty log file is a valid state.
func TestTailEmptyFile(t *testing.T) {
	path := writeTempLog(t, "")
	r, err := tail.New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	lines, err := r.Lines(10)
	if err != nil {
		t.Fatalf("Lines on empty file: %v", err)
	}
	if len(lines) != 0 {
		t.Errorf("expected no lines for empty file, got %d: %v", len(lines), lines)
	}
}
