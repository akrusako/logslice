package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/yourorg/logslice/internal/ratelimit"

	"github.com/spf13/cobra"
)

func newRateLimitCmd() *cobra.Command {
	var maxPerSec int
	var maxTotal int

	cmd := &cobra.Command{
		Use:   "ratelimit [file]",
		Short: "Stream a log file with output throttling",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var f *os.File
			if len(args) == 0 {
				f = os.Stdin
			} else {
				var err error
				f, err = os.Open(args[0])
				if err != nil {
					return fmt.Errorf("open: %w", err)
				}
				defer f.Close()
			}

			limiter := ratelimit.New(maxPerSec, maxTotal)
			scanner := bufio.NewScanner(f)
			w := bufio.NewWriter(os.Stdout)
			defer w.Flush()

			for scanner.Scan() {
				if !limiter.Allow() {
					break
				}
				fmt.Fprintln(w, scanner.Text())
			}
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("scan: %w", err)
			}
			return nil
		},
	}

	cmd.Flags().IntVar(&maxPerSec, "rate", 0, "max lines per second (0 = unlimited)")
	cmd.Flags().IntVar(&maxTotal, "max", 0, "max total lines to emit (0 = unlimited)")
	return cmd
}
