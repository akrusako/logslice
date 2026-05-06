package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/user/logslice/internal/sample"
	"github.com/spf13/cobra"
)

func newSampleCmd() *cobra.Command {
	var rate float64
	var seed int64

	cmd := &cobra.Command{
		Use:   "sample <file>",
		Short: "Probabilistically sample lines from a log file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			f, err := os.Open(args[0])
			if err != nil {
				return fmt.Errorf("open: %w", err)
			}
			defer f.Close()

			if seed == 0 {
				seed = time.Now().UnixNano()
			}

			s := sample.New(rate, seed)
			sc := bufio.NewScanner(f)
			w := bufio.NewWriter(os.Stdout)
			defer w.Flush()

			for sc.Scan() {
				line := sc.Text()
				if s.Keep(line) {
					fmt.Fprintln(w, line)
				}
			}
			if err := sc.Err(); err != nil {
				return fmt.Errorf("scan: %w", err)
			}

			total, kept := s.Stats()
			fmt.Fprintf(os.Stderr, "sampled %d/%d lines (rate=%.2f)\n", kept, total, rate)
			return nil
		},
	}

	cmd.Flags().Float64VarP(&rate, "rate", "r", 1.0, "sampling rate between 0.0 and 1.0")
	cmd.Flags().Int64Var(&seed, "seed", 0, "random seed (0 = use current time)")
	return cmd
}
