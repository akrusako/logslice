package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempRedactLog(t *testing.T, lines []string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "redact*.log")
	if err != nil {
		t.Fatalf("temp file: %v", err)
	}
	for _, l := range lines {
		f.WriteString(l + "\n")
	}
	f.Close()
	return f.Name()
}

func TestRedactCmdNoArgs(t *testing.T) {
	cmd := newRedactCmd()
	cmd.SetArgs([]string{})
	if err := cmd.Execute(); err == nil {
		t.Error("expected error with no args")
	}
}

func TestRedactCmdBadFile(t *testing.T) {
	cmd := newRedactCmd()
	cmd.SetArgs([]string{"--pattern", `\d+`, filepath.Join(t.TempDir(), "missing.log")})
	if err := cmd.Execute(); err == nil {
		t.Error("expected error for missing file")
	}
}

func TestRedactCmdNoPattern(t *testing.T) {
	f := writeTempRedactLog(t, []string{"hello"})
	cmd := newRedactCmd()
	cmd.SetArgs([]string{f})
	if err := cmd.Execute(); err == nil {
		t.Error("expected error when no pattern given")
	}
}

func TestRedactCmdReplacesPattern(t *testing.T) {
	lines := []string{
		"user=alice@example.com connected",
		"no sensitive data here",
		"card=4111-1111-1111-1111 charged",
	}
	f := writeTempRedactLog(t, lines)

	var sb strings.Builder
	cmd := newRedactCmd()
	cmd.SetOut(&sb)
	cmd.SetArgs([]string{
		"--pattern", `[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}`,
		"--pattern", `\d{4}-\d{4}-\d{4}-\d{4}`,
		"--mask", "[X]",
		"--quiet",
		f,
	})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if strings.Contains(out, "alice@example.com") {
		t.Error("email not redacted")
	}
	if strings.Contains(out, "4111-1111-1111-1111") {
		t.Error("card not redacted")
	}
	if !strings.Contains(out, "no sensitive data here") {
		t.Error("plain line should pass through unchanged")
	}
}
