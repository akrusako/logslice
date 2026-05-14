package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/user/logslice/internal/fieldsorter"
	"github.com/spf13/cobra"
)

func newFieldSorterCmd() *cobra.Command {
	var fields string

	cmd := &cobra.Command{
		Use:   "fieldsorter <file>",
		Short: "Reorder fields in structured log lines",
		Long: `Reads a log file and reorders fields in each line so that
the specified fields appear first, in the given order.
Supports JSON and key=value formats.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if fields == "" {
				return fmt.Errorf("--fields is required")
			}
			order := strings.Split(fields, ",")
			for i, f := range order {
				order[i] = strings.TrimSpace(f)
			}

			s, err := fieldsorter.New(order)
			if err != nil {
				return err
			}

			f, err := os.Open(args[0])
			if err != nil {
				return fmt.Errorf("open %s: %w", args[0], err)
			}
			defer f.Close()

			scanner := bufio.NewScanner(f)
			w := bufio.NewWriter(os.Stdout)
			defer w.Flush()

			for scanner.Scan() {
				fmt.Fprintln(w, s.Apply(scanner.Text()))
			}
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("scan: %w", err)
			}

			fmt.Fprintf(os.Stderr, "sorted=%d skipped=%d\n",
				s.Sorted(), s.Skipped())
			return nil
		},
	}

	cmd.Flags().StringVar(&fields, "fields", "",
		"comma-separated list of fields to place first (e.g. level,msg,ts)")
	return cmd
}
