package linecount_test

import (
	"os"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/linecount"
)

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "linecount-*.log")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestReaderEmptyInput(t *testing.T) {
	res, err := linecount.Reader(strings.NewReader(""))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Lines != 0 || res.Bytes != 0 {
		t.Errorf("expected 0 lines/bytes, got %+v", res)
	}
}

func TestReaderNewlineTerminated(t *testing.T) {
	input := "line1\nline2\nline3\n"
	res, err := linecount.Reader(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Lines != 3 {
		t.Errorf("expected 3 lines, got %d", res.Lines)
	}
	if res.Bytes != int64(len(input)) {
		t.Errorf("expected %d bytes, got %d", len(input), res.Bytes)
	}
}

func TestReaderNoTrailingNewline(t *testing.T) {
	input := "alpha\nbeta\ngamma"
	res, err := linecount.Reader(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Lines != 3 {
		t.Errorf("expected 3 lines, got %d", res.Lines)
	}
}

func TestFileCounting(t *testing.T) {
	path := writeTempFile(t, "a\nb\nc\nd\n")
	res, err := linecount.File(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Lines != 4 {
		t.Errorf("expected 4 lines, got %d", res.Lines)
	}
}

func TestFileBadPath(t *testing.T) {
	_, err := linecount.File("/nonexistent/path/file.log")
	if err == nil {
		t.Error("expected error for bad path, got nil")
	}
}

func TestFromOffset(t *testing.T) {
	// First two lines: "foo\nbar\n" = 8 bytes; remaining: "baz\n"
	path := writeTempFile(t, "foo\nbar\nbaz\n")
	res, err := linecount.FromOffset(path, 8)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Lines != 1 {
		t.Errorf("expected 1 line from offset, got %d", res.Lines)
	}
}
