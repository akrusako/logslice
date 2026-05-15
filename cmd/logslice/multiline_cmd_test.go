package main

import (
	"os"
	"strings"
	"testing"
)

func writeTempMultilineLog(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "multiline-*.log")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestMultilineCmdNoArgs(t *testing.T) {
	cmd := newMultilineCmd()
	cmd.SetArgs([]string{})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error with no args")
	}
}

func TestMultilineCmdBadFile(t *testing.T) {
	cmd := newMultilineCmd()
	cmd.SetArgs([]string{"--pattern", `^\d`, "/nonexistent/file.log"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestMultilineCmdMissingPattern(t *testing.T) {
	path := writeTempMultilineLog(t, "2024-01-01 hello\n")
	cmd := newMultilineCmd()
	cmd.SetArgs([]string{path})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error when --pattern is missing")
	}
}

func TestMultilineCmdMergesLines(t *testing.T) {
	content := "2024-01-01 start\n  cont one\n  cont two\n2024-01-02 next\n"
	path := writeTempMultilineLog(t, content)

	var sb strings.Builder
	cmd := newMultilineCmd()
	cmd.SetOut(&sb)
	cmd.SetArgs([]string{"--pattern", `^\d{4}-`, "--sep", " ", path})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := sb.String()
	if !strings.Contains(output, "cont one") {
		t.Errorf("expected merged continuation in output, got: %q", output)
	}
	if strings.Count(output, "\n") < 2 {
		t.Errorf("expected at least 2 output lines, got: %q", output)
	}
}
