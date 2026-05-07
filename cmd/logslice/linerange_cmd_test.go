package main

import (
	"os"
	"strings"
	"testing"
)

func writeTempLineLog(t *testing.T, lines []string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "linerange-*.log")
	if err != nil {
		t.Fatal(err)
	}
	for _, l := range lines {
		f.WriteString(l + "\n")
	}
	f.Close()
	return f.Name()
}

func TestLineRangeCmdNoArgs(t *testing.T) {
	cmd := newLineRangeCmd()
	cmd.SetArgs([]string{})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error with no arguments")
	}
}

func TestLineRangeCmdBadFile(t *testing.T) {
	cmd := newLineRangeCmd()
	cmd.SetArgs([]string{"--range", "1:5", "/nonexistent/file.log"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLineRangeCmdMissingRange(t *testing.T) {
	lines := []string{"a", "b", "c"}
	path := writeTempLineLog(t, lines)
	cmd := newLineRangeCmd()
	cmd.SetArgs([]string{path})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error when --range not provided")
	}
}

func TestLineRangeCmdSubset(t *testing.T) {
	lines := []string{"line1", "line2", "line3", "line4", "line5"}
	path := writeTempLineLog(t, lines)

	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w

	cmd := newLineRangeCmd()
	cmd.SetArgs([]string{"--range", "2:4", path})
	if err := cmd.Execute(); err != nil {
		w.Close()
		os.Stdout = old
		t.Fatalf("unexpected error: %v", err)
	}
	w.Close()
	os.Stdout = old

	buf := new(strings.Builder)
	buf2 := make([]byte, 1024)
	n, _ := r.Read(buf2)
	buf.Write(buf2[:n])

	got := strings.TrimSpace(buf.String())
	want := "line2\nline3\nline4"
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}
