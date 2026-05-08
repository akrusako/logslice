package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func writeTempWindowLog(t *testing.T, lines []string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "window*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	for _, l := range lines {
		f.WriteString(l + "\n")
	}
	return f.Name()
}

func TestWindowCmdNoArgs(t *testing.T) {
	cmd := newWindowCmd()
	cmd.SetArgs([]string{})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error with no args")
	}
}

func TestWindowCmdBadFile(t *testing.T) {
	cmd := newWindowCmd()
	cmd.SetArgs([]string{"/nonexistent/file.log"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestWindowCmdCountsFlag(t *testing.T) {
	lines := []string{
		"2024-01-01T00:00:01Z level=info msg=a",
		"2024-01-01T00:00:02Z level=info msg=b",
		"2024-01-01T00:00:03Z level=info msg=c",
	}
	path := writeTempWindowLog(t, lines)

	var out bytes.Buffer
	cmd := newWindowCmd()
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"--counts", "--window", "60s", path})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result := out.String()
	if !strings.Contains(result, "1\t") {
		t.Errorf("expected count prefix in output, got: %q", result)
	}
}

func TestWindowCmdThreshold(t *testing.T) {
	lines := []string{
		"2024-01-01T00:00:01Z msg=one",
		"2024-01-01T00:00:02Z msg=two",
		"2024-01-01T00:00:03Z msg=three",
		"2024-01-01T00:00:04Z msg=four",
	}
	path := writeTempWindowLog(t, lines)

	var out bytes.Buffer
	cmd := newWindowCmd()
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"--threshold", "2", "--window", "60s", path})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result := out.String()
	// lines 3 and 4 have window counts 3 and 4 (> threshold 2)
	if !strings.Contains(result, "three") {
		t.Errorf("expected 'three' in output, got: %q", result)
	}
	if strings.Contains(result, "one") {
		t.Errorf("did not expect 'one' in output, got: %q", result)
	}
}
