package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/logslice/internal/columnar"
)

func newColumnarCmd() *cobra.Command {
	var (
		fields  string
		width   int
		sep     string
		header  bool
	)

	cmd := &cobra.Command{
		Use:   "columnar <file>",
		Short: "Extract log fields as aligned columns",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if fields == "" {
				return fmt.Errorf("--fields is required")
			}

			f, err := os.Open(args[0])
			if err != nil {
				return fmt.Errorf("open: %w", err)
			}
			defer f.Close()

			names := strings.Split(fields, ",")
			for i, n := range names {
				names[i] = strings.TrimSpace(n)
			}

			c := columnar.New(names, width, sep)

			if header {
				fmt.Fprintln(cmd.OutOrStdout(), c.Header())
			}

			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				row, _ := c.Format(scanner.Text())
				fmt.Fprintln(cmd.OutOrStdout(), row)
			}
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("scan: %w", err)
			}

			if d := c.DroppedCount(); d > 0 {
				fmt.Fprintf(cmd.ErrOrStderr(), "warn: %d line(s) had missing fields\n", d)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&fields, "fields", "", "Comma-separated list of fields to extract (required)")
	cmd.Flags().IntVar(&width, "width", 0, "Minimum column width for padding (0=disabled)")
	cmd.Flags().StringVar(&sep, "sep", "\t", "Column separator")
	cmd.Flags().BoolVar(&header, "header", false, "Print a header row")

	return cmd
}
