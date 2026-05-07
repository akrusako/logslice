package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/yourorg/logslice/internal/throttle"
	"github.com/spf13/cobra"
)

func newThrottleCmd() *cobra.Command {
	var (
		window   time.Duration
		maxLines int
	)

	cmd := &cobra.Command{
		Use:   "throttle <file>",
		Short: "Limit log output to N lines per time window",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			f, err := os.Open(args[0])
			if err != nil {
				return fmt.Errorf("open: %w", err)
			}
			defer f.Close()

			th := throttle.New(window, maxLines)
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				if th.Allow() {
					fmt.Fprintln(cmd.OutOrStdout(), scanner.Text())
				}
			}
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("scan: %w", err)
			}

			fmt.Fprintf(
				cmd.ErrOrStderr(),
				"throttle: %d/%d lines emitted (%d dropped)\n",
				th.Total()-th.Dropped(), th.Total(), th.Dropped(),
			)
			return nil
		},
	}

	cmd.Flags().DurationVarP(&window, "window", "w", time.Second,
		"sliding time window for the rate cap")
	cmd.Flags().IntVarP(&maxLines, "max", "n", 100,
		"maximum lines allowed per window")
	return cmd
}
