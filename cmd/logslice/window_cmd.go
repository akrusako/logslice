package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/yourorg/logslice/internal/timeparse"
	"github.com/yourorg/logslice/internal/window"
	"github.com/spf13/cobra"
)

func newWindowCmd() *cobra.Command {
	var (
		windowSize  time.Duration
		threshold   int64
		printCounts bool
	)

	cmd := &cobra.Command{
		Use:   "window <file>",
		Short: "Count lines per sliding time window and flag spikes",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			f, err := os.Open(args[0])
			if err != nil {
				return fmt.Errorf("open: %w", err)
			}
			defer f.Close()

			agg := window.New(windowSize)
			parser := timeparse.NewParser()
			scanner := bufio.NewScanner(f)

			for scanner.Scan() {
				line := scanner.Text()
				t, err := parser.ParseAny(line)
				if err != nil {
					t = time.Now()
				}
				count := agg.Add(t)
				if printCounts {
					fmt.Fprintf(cmd.OutOrStdout(), "%d\t%s\n", count, line)
					continue
				}
				if threshold > 0 && count > threshold {
					fmt.Fprintln(cmd.OutOrStdout(), line)
				}
			}
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("scan: %w", err)
			}
			fmt.Fprintf(cmd.ErrOrStderr(), "total lines: %d\n", agg.Total())
			return nil
		},
	}

	cmd.Flags().DurationVarP(&windowSize, "window", "w", 60*time.Second, "sliding window size")
	cmd.Flags().Int64VarP(&threshold, "threshold", "t", 0, "emit lines when window count exceeds this value (0 = disabled)")
	cmd.Flags().BoolVar(&printCounts, "counts", false, "prefix each line with its window count")
	return cmd
}
