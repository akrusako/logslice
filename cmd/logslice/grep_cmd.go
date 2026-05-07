package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/yourorg/logslice/internal/grep"
	"github.com/spf13/cobra"
)

func newGrepCmd() *cobra.Command {
	var invert bool

	cmd := &cobra.Command{
		Use:   "grep [flags] <file> <pattern> [pattern...]",
		Short: "Filter log lines by regex pattern",
		Long: `grep reads a log file and prints only lines matching one or more
regex patterns. Use --invert to exclude matching lines instead.`,
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]
			patterns := args[1:]

			m, err := grep.New(patterns, invert)
			if err != nil {
				return err
			}

			f, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("grep: open %q: %w", path, err)
			}
			defer f.Close()

			sc := bufio.NewScanner(f)
			w := bufio.NewWriter(os.Stdout)
			defer w.Flush()

			for sc.Scan() {
				line := sc.Text()
				if m.Match(line) {
					fmt.Fprintln(w, line)
				}
			}
			if err := sc.Err(); err != nil {
				return fmt.Errorf("grep: scan: %w", err)
			}

			fmt.Fprintf(os.Stderr, "grep: matched=%d dropped=%d\n",
				m.Matched(), m.Dropped())
			return nil
		},
	}

	cmd.Flags().BoolVarP(&invert, "invert", "v", false,
		"invert match: print lines that do NOT match any pattern")
	return cmd
}
