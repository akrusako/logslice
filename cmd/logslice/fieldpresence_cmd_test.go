package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempPresenceLog(t *testing.T, lines []string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "presence*.log")
	if err != nil {
		t.Fatal(err)
	}
	for _, l := range lines {
		f.WriteString(l + "\n")
	}
	f.Close()
	return filepath.Clean(f.Name())
}

func TestFieldPresenceCmdNoArgs(t *testing.T) {
	cmd := newFieldPresenceCmd()
	cmd.SetArgs([]string{})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error with no args")
	}
}

func TestFieldPresenceCmdBadFile(t *testing.T) {
	cmd := newFieldPresenceCmd()
	cmd.SetArgs([]string{"--require", "level", "/nonexistent/file.log"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestFieldPresenceCmdNoFlags(t *testing.T) {
	p := writeTempPresenceLog(t, []string{"level=info"})
	cmd := newFieldPresenceCmd()
	cmd.SetArgs([]string{p})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error when neither --require nor --forbid set")
	}
}

func TestFieldPresenceCmdRequireFilters(t *testing.T) {
	lines := []string{
		"level=info msg=hello",
		"msg=no_level",
		`{"level":"warn","msg":"json"}`,
	}
	p := writeTempPresenceLog(t, lines)

	var out bytes.Buffer
	cmd := newFieldPresenceCmd()
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"--require", "level", p})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result := out.String()
	if !strings.Contains(result, "level=info") {
		t.Error("expected kv line with level to appear")
	}
	if strings.Contains(result, "msg=no_level") {
		t.Error("expected line without level to be dropped")
	}
	if !strings.Contains(result, "warn") {
		t.Error("expected JSON line with level to appear")
	}
}
