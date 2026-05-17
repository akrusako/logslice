package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/yourorg/logslice/internal/fieldformat"
	"github.com/spf13/cobra"
)

func newFieldFormatCmd() *cobra.Command {
	var field, format string

	cmd := &cobra.Command{
		Use:   "fieldformat <file>",
		Short: "Reformat a log field value using a Go format string",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if field == "" {
				return fmt.Errorf("--field is required")
			}
			if format == "" {
				return fmt.Errorf("--format is required")
			}

			fmt_, err := fieldformat.New(field, format)
			if err != nil {
				return err
			}

			f, err := os.Open(args[0])
			if err != nil {
				return fmt.Errorf("fieldformat: %w", err)
			}
			defer f.Close()

			w := bufio.NewWriter(os.Stdout)
			defer w.Flush()

			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				fmt.Fprintln(w, fmt_.Apply(scanner.Text()))
			}
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("fieldformat: read error: %w", err)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&field, "field", "", "field name to reformat")
	cmd.Flags().StringVar(&format, "format", "", "Go format string (e.g. \"%.2f\")")
	return cmd
}
