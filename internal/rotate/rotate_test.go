package rotate

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "test.log")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	return p
}

func TestNewOpensFile(t *testing.T) {
	p := writeTempFile(t, "hello\n")
	w, err := New(p, 100*time.Millisecond)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer w.Close()
	if w.Reader() == nil {
		t.Fatal("expected non-nil reader")
	}
}

func TestNewBadPath(t *testing.T) {
	_, err := New("/nonexistent/path/file.log", 0)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestNotRotatedWhenUnchanged(t *testing.T) {
	p := writeTempFile(t, "line1\n")
	w, err := New(p, 100*time.Millisecond)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer w.Close()

	rotated, err := w.Rotated()
	if err != nil {
		t.Fatalf("Rotated: %v", err)
	}
	if rotated {
		t.Fatal("expected not rotated for unchanged file")
	}
}

func TestRotatedOnTruncation(t *testing.T) {
	p := writeTempFile(t, "line1\nline2\nline3\n")
	w, err := New(p, 100*time.Millisecond)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer w.Close()

	// Truncate the file to simulate rotation.
	if err := os.WriteFile(p, []byte(""), 0o644); err != nil {
		t.Fatalf("truncate: %v", err)
	}

	rotated, err := w.Rotated()
	if err != nil {
		t.Fatalf("Rotated: %v", err)
	}
	if !rotated {
		t.Fatal("expected rotation detected after truncation")
	}
}

func TestReopenResetsState(t *testing.T) {
	p := writeTempFile(t, "data\n")
	w, err := New(p, 100*time.Millisecond)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer w.Close()

	if err := w.Reopen(); err != nil {
		t.Fatalf("Reopen: %v", err)
	}
	if w.Reader() == nil {
		t.Fatal("expected non-nil reader after reopen")
	}
}

func TestPollIntervalDefault(t *testing.T) {
	p := writeTempFile(t, "x\n")
	w, err := New(p, 0)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer w.Close()
	if w.PollInterval() != 500*time.Millisecond {
		t.Fatalf("expected default 500ms, got %v", w.PollInterval())
	}
}
