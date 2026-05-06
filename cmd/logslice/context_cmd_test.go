package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempContextLog(t *testing.T, lines []string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "ctx-*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	for _, l := range lines {
		f.WriteString(l + "\n")
	}
	return f.Name()
}

func TestContextCmdNoArgs(t *testing.T) {
	cmd := newContextCmd()
	cmd.SetArgs([]string{})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error with no args")
	}
}

func TestContextCmdBadFile(t *testing.T) {
	cmd := newContextCmd()
	cmd.SetArgs([]string{filepath.Join(t.TempDir(), "missing.log")})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestContextCmdMaxLines(t *testing.T) {
	lines := []string{"alpha", "beta", "gamma", "delta", "epsilon"}
	path := writeTempContextLog(t, lines)

	var buf strings.Builder
	cmd := newContextCmd()
	cmd.SetArgs([]string{"--max-lines=3", "--quiet", path})
	cmd.SetOut(&buf)
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	got := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(got) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(got))
	}
}

func TestContextCmdBadSince(t *testing.T) {
	path := writeTempContextLog(t, []string{"line"})
	cmd := newContextCmd()
	cmd.SetArgs([]string{"--since=not-a-date", path})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for invalid --since")
	}
}
