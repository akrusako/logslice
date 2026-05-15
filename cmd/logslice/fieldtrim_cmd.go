package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/yourorg/logslice/internal/fieldtrim"
	"github.com/spf13/cobra"
)

func newFieldTrimCmd() *cobra.Command {
	var cutset string

	cmd := &cobra.Command{
		Use:   "fieldtrim <file>",
		Short: "Trim whitespace or a custom cutset from a named field value",
		Example: `  logslice fieldtrim --field msg app.log
  logslice fieldtrim --field path --cutset / app.log`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			field, _ := cmd.Flags().GetString("field")
			if field == "" {
				return fmt.Errorf("--field is required")
			}

			tr, err := fieldtrim.New(field, cutset)
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
				fmt.Fprintln(w, tr.Apply(scanner.Text()))
			}
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("scan: %w", err)
			}
			fmt.Fprintf(os.Stderr, "fieldtrim: %d line(s) trimmed\n", tr.Trimmed())
			return nil
		},
	}

	cmd.Flags().String("field", "", "field name whose value should be trimmed (required)")
	cmd.Flags().StringVar(&cutset, "cutset", "", "characters to remove from both ends (default: Unicode whitespace)")
	return cmd
}
