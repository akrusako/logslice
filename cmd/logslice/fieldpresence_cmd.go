package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/logslice/internal/fieldpresence"
)

func newFieldPresenceCmd() *cobra.Command {
	var require []string
	var forbid []string

	cmd := &cobra.Command{
		Use:   "fieldpresence <file>",
		Short: "Filter lines by field presence or absence",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(require) == 0 && len(forbid) == 0 {
				return fmt.Errorf("at least one --require or --forbid flag must be set")
			}

			checker, err := fieldpresence.New(require, forbid)
			if err != nil {
				return err
			}

			f, err := os.Open(args[0])
			if err != nil {
				return fmt.Errorf("open %s: %w", args[0], err)
			}
			defer f.Close()

			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				line := scanner.Text()
				if checker.Allow(line) {
					fmt.Fprintln(cmd.OutOrStdout(), line)
				}
			}
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("scan: %w", err)
			}

			fmt.Fprintf(cmd.ErrOrStderr(), "checked=%d passed=%d\n",
				checker.Checked(), checker.Passed())
			return nil
		},
	}

	cmd.Flags().StringSliceVar(&require, "require", nil,
		"comma-separated fields that must be present (e.g. level,request_id)")
	cmd.Flags().StringSliceVar(&forbid, "forbid", nil,
		"comma-separated fields that must be absent (e.g. debug,trace)")
	_ = strings.Join // keep import used via StringSliceVar internals
	return cmd
}
