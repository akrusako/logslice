package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempRateLog(t *testing.T, lines int) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "rate.log")
	f, err := os.Create(p)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	for i := 0; i < lines; i++ {
		fmt.Fprintf(f, "line %d\n", i+1)
	}
	return p
}

func TestRateLimitCmdNoArgs(t *testing.T) {
	cmd := newRateLimitCmd()
	cmd.SetArgs([]string{"--max", "5", "--"})
	// passing no file and no stdin redirect; just ensure command is wired
	if cmd == nil {
		t.Fatal("expected non-nil command")
	}
}

func TestRateLimitCmdBadFile(t *testing.T) {
	cmd := newRateLimitCmd()
	cmd.SetArgs([]string{"/no/such/file.log"})
	var buf bytes.Buffer
	cmd.SetErr(&buf)
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestRateLimitCmdMaxTotal(t *testing.T) {
	p := writeTempRateLog(t, 20)
	cmd := newRateLimitCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"--max", "7", p})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(got) != 7 {
		t.Fatalf("expected 7 lines, got %d", len(got))
	}
}
