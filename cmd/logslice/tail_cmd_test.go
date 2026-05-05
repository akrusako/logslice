package main

import (
	"flag"
	"os"
	"testing"
)

func TestTailCmdNoArgs(t *testing.T) {
	cmd := &tailCmd{}
	fs := flag.NewFlagSet("tail", flag.ContinueOnError)
	cmd.register(fs)

	err := cmd.run(nil)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestTailCmdBadFile(t *testing.T) {
	cmd := &tailCmd{n: 5, fmt: "raw"}
	fs := flag.NewFlagSet("tail", flag.ContinueOnError)
	cmd.register(fs)

	err := cmd.run([]string{"/nonexistent/path/file.log"})
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestTailCmdRuns(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "tail-cmd-*.log")
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 20; i++ {
		fmt.Fprintf(f, "2024-01-01T00:00:%02dZ level=info msg=line%d\n", i, i)
	}
	f.Close()

	cmd := &tailCmd{n: 5, fmt: "raw"}
	fs := flag.NewFlagSet("tail", flag.ContinueOnError)
	cmd.register(fs)

	if err := cmd.run([]string{f.Name()}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
