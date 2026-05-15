package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/user/logslice/internal/fielddefault"
	"github.com/spf13/cobra"
)

func newFieldDefaultCmd() *cobra.Command {
	var field string
	var value string

	cmd := &cobra.Command{
		Use:   "fielddefault <file>",
		Short: "Inject a default value for a missing field in each log line",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if field == "" {
				return fmt.Errorf("--field is required")
			}
			d, err := fielddefault.New(field, value)
			if err != nil {
				return err
			}
			f, err := os.Open(args[0])
			if err != nil {
				return err
			}
			defer f.Close()

			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				fmt.Println(d.Apply(scanner.Text()))
			}
			if err := scanner.Err(); err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "applied=%d skipped=%d\n", d.Applied(), d.Skipped())
			return nil
		},
	}

	cmd.Flags().StringVar(&field, "field", "", "Field name to check and inject")
	cmd.Flags().StringVar(&value, "value", "", "Default value to inject when field is absent")
	return cmd
}
