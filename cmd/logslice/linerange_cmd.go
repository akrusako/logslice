package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/yourorg/logslice/internal/linerange"
	"github.com/spf13/cobra"
)

func newLineRangeCmd() *cobra.Command {
	var rangeStr string

	cmd := &cobra.Command{
		Use:   "linerange <file>",
		Short: "Extract lines by line-number range from a log file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if rangeStr == "" {
				return fmt.Errorf("--range is required")
			}

			r, err := linerange.Parse(rangeStr)
			if err != nil {
				return fmt.Errorf("invalid range: %w", err)
			}

			f, err := os.Open(args[0])
			if err != nil {
				return fmt.Errorf("open: %w", err)
			}
			defer f.Close()

			out := bufio.NewWriter(os.Stdout)
			defer out.Flush()

			scanner := bufio.NewScanner(f)
			lineNum := 0
			for scanner.Scan() {
				lineNum++
				// If we have a bounded range and we've passed it, stop early.
				if r.End != 0 && lineNum > r.End {
					break
				}
				if r.Contains(lineNum) {
					fmt.Fprintln(out, scanner.Text())
				}
			}
			return scanner.Err()
		},
	}

	cmd.Flags().StringVar(&rangeStr, "range", "", `line range to extract, e.g. "10:50", "100:", ":25"`)
	return cmd
}
