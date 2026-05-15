package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/yourorg/logslice/internal/multiline"
	"github.com/spf13/cobra"
)

func newMultilineCmd() *cobra.Command {
	var (
		pattern string
		sep     string
	)

	cmd := &cobra.Command{
		Use:   "multiline <file>",
		Short: "Merge continuation lines into single log entries",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if pattern == "" {
				return fmt.Errorf("--pattern is required")
			}
			f, err := os.Open(args[0])
			if err != nil {
				return fmt.Errorf("multiline: open file: %w", err)
			}
			defer f.Close()

			merger, err := multiline.New(pattern, sep)
			if err != nil {
				return err
			}

			out := bufio.NewWriter(os.Stdout)
			defer out.Flush()

			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				if entry, ok := merger.Feed(scanner.Text()); ok {
					fmt.Fprintln(out, entry)
				}
			}
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("multiline: scan: %w", err)
			}
			if entry, ok := merger.Flush(); ok {
				fmt.Fprintln(out, entry)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&pattern, "pattern", "", "Regex that marks the start of a new log entry")
	cmd.Flags().StringVar(&sep, "sep", " ", "Separator used when joining continuation lines")
	return cmd
}
