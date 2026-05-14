package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempSorterLog(t *testing.T, lines []string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "sorterlog-*.log")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	defer f.Close()
	for _, l := range lines {
		f.WriteString(l + "\n")
	}
	return filepath.Clean(f.Name())
}

func TestFieldSorterCmdNoArgs(t *testing.T) {
	cmd := newFieldSorterCmd()
	cmd.SetArgs([]string{})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error with no args")
	}
}

func TestFieldSorterCmdBadFile(t *testing.T) {
	cmd := newFieldSorterCmd()
	cmd.SetArgs([]string{"--fields", "level,msg", "/nonexistent/file.log"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for bad file")
	}
}

func TestFieldSorterCmdMissingFields(t *testing.T) {
	path := writeTempSorterLog(t, []string{"level=info msg=hello"})
	cmd := newFieldSorterCmd()
	cmd.SetArgs([]string{path})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error when --fields not provided")
	}
}

func TestFieldSorterCmdRuns(t *testing.T) {
	lines := []string{
		"msg=hello level=info ts=2024-01-01",
		`{"ts":"2024-01-01","msg":"bye","level":"warn"}`,
	}
	path := writeTempSorterLog(t, lines)

	r, w, _ := os.Pipe()
	oldStdout := os.Stdout
	os.Stdout = w

	cmd := newFieldSorterCmd()
	cmd.SetArgs([]string{"--fields", "level,msg", path})
	err := cmd.Execute()

	w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	buf := new(strings.Builder)
	buf2 := make([]byte, 4096)
	n, _ := r.Read(buf2)
	buf.Write(buf2[:n])
	out := buf.String()

	if !strings.Contains(out, "level=info msg=hello") {
		t.Errorf("expected sorted KV line in output, got:\n%s", out)
	}
}
