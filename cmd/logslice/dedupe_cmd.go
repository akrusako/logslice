package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/yourorg/logslice/internal/dedupe"
	"github.com/spf13/cobra"
)

func newDedupeCmd() *cobra.Command {
	var windowSize int

	cmd := &cobra.Command{
		Use:   "dedupe <file>",
		Short: "Remove duplicate log lines within a sliding window",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]
			f, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("open %s: %w", path, err)
			}
			defer f.Close()

			d := dedupe.New(windowSize)
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				line := scanner.Text()
				if !d.IsDuplicate(line) {
					fmt.Fprintln(cmd.OutOrStdout(), line)
				}
			}
			return scanner.Err()
		},
	}

	cmd.Flags().IntVarP(&windowSize, "window", "w", 1000,
		"number of recent lines to track for deduplication (0 = disabled)")
	// expose window size in help as human-readable string
	_ = strconv.Itoa(windowSize)
	return cmd
}
