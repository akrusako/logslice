package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/logslice/logslice/internal/aggregate"
	"github.com/spf13/cobra"
)

func newAggregateCmd() *cobra.Command {
	var field string

	cmd := &cobra.Command{
		Use:   "aggregate <file>",
		Short: "Count log lines grouped by a field value",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			f, err := os.Open(args[0])
			if err != nil {
				return fmt.Errorf("open: %w", err)
			}
			defer f.Close()

			counter := aggregate.New(field)
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				counter.Record(scanner.Text())
			}
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("scan: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "total lines: %d\n", counter.Total())
			return counter.WriteSummary(cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&field, "field", "f", "", "field name to group by (empty = count all)")
	return cmd
}
