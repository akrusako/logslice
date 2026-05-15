package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempCastLog(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "cast.log")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempCastLog: %v", err)
	}
	return p
}

func TestFieldCastCmdNoArgs(t *testing.T) {
	cmd := newFieldCastCmd()
	cmd.SetArgs([]string{})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error with no args")
	}
}

func TestFieldCastCmdBadFile(t *testing.T) {
	cmd := newFieldCastCmd()
	cmd.SetArgs([]string{"/nonexistent/file.log", "latency", "--type", "int"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestFieldCastCmdBadType(t *testing.T) {
	p := writeTempCastLog(t, "latency=3.2 msg=ok\n")
	cmd := newFieldCastCmd()
	cmd.SetArgs([]string{p, "latency", "--type", "bytes"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for unknown type")
	}
}

func TestFieldCastCmdIntConversion(t *testing.T) {
	lines := "latency=3.9 msg=ok\nlatency=1.1 msg=hi\n"
	p := writeTempCastLog(t, lines)

	buf := &strings.Builder{}
	cmd := newFieldCastCmd()
	cmd.SetOut(buf)
	cmd.SetArgs([]string{p, "latency", "--type", "int"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if strings.Contains(out, "3.9") || strings.Contains(out, "1.1") {
		t.Errorf("expected int values in output, got: %q", out)
	}
}
