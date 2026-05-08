package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/yourorg/logslice/internal/redact"
	"github.com/spf13/cobra"
)

func newRedactCmd() *cobra.Command {
	var (
		patterns []string
		mask     string
		quiet    bool
	)

	cmd := &cobra.Command{
		Use:   "redact <file>",
		Short: "Replace sensitive patterns in log lines",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(patterns) == 0 {
				return fmt.Errorf("at least one --pattern is required")
			}
			r, err := redact.New(patterns, mask)
			if err != nil {
				return err
			}
			f, err := os.Open(args[0])
			if err != nil {
				return err
			}
			defer f.Close()

			scanner := bufio.NewScanner(f)
			w := bufio.NewWriter(os.Stdout)
			defer w.Flush()

			for scanner.Scan() {
				out, _ := r.Apply(scanner.Text())
				fmt.Fprintln(w, out)
			}
			if err := scanner.Err(); err != nil {
				return err
			}
			if !quiet {
				fmt.Fprintf(os.Stderr, "redacted lines: %d\n", r.Count())
			}
			return nil
		},
	}

	cmd.Flags().StringArrayVarP(&patterns, "pattern", "p", nil, "regex pattern to redact (repeatable)")
	cmd.Flags().StringVarP(&mask, "mask", "m", "", "replacement string (default [REDACTED])")
	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "suppress summary output")
	return cmd
}
