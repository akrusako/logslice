package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func writeTempDedupLog(t *testing.T, lines []string) string {
	t.Helper()
	f, err := os.CreateTemp("", "dedupe-*.log")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	t.Cleanup(func() { os.Remove(f.Name()) })
	for _, l := range lines {
		f.WriteString(l + "\n")
	}
	f.Close()
	return f.Name()
}

func TestDedupeCmdNoArgs(t *testing.T) {
	cmd := newDedupeCmd()
	cmd.SetArgs([]string{})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error with no arguments")
	}
}

func TestDedupeCmdBadFile(t *testing.T) {
	cmd := newDedupeCmd()
	cmd.SetArgs([]string{"/nonexistent/path/file.log"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestDedupeCmdRemovesDuplicates(t *testing.T) {
	lines := []string{"alpha", "beta", "alpha", "gamma", "beta", "delta"}
	path := writeTempDedupLog(t, lines)

	var buf bytes.Buffer
	cmd := newDedupeCmd()
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{path})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := strings.Split(strings.TrimSpace(buf.String()), "\n")
	want := []string{"alpha", "beta", "gamma", "delta"}
	if len(got) != len(want) {
		t.Fatalf("got %d lines, want %d: %v", len(got), len(want), got)
	}
	for i, w := range want {
		if got[i] != w {
			t.Errorf("line %d: got %q, want %q", i, got[i], w)
		}
	}
}
