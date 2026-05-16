package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/yourorg/logslice/internal/fieldcopy"
	"github.com/spf13/cobra"
)

func newFieldCopyCmd() *cobra.Command {
	var src, dst string

	cmd := &cobra.Command{
		Use:   "fieldcopy <file>",
		Short: "Duplicate a field value under a new key in each log line",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if src == "" {
				return fmt.Errorf("--src is required")
			}
			if dst == "" {
				return fmt.Errorf("--dst is required")
			}

			copier, err := fieldcopy.New(src, dst)
			if err != nil {
				return err
			}

			f, err := os.Open(args[0])
			if err != nil {
				return fmt.Errorf("fieldcopy: open file: %w", err)
			}
			defer f.Close()

			w := bufio.NewWriter(os.Stdout)
			defer w.Flush()

			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				line := scanner.Text()
				fmt.Fprintln(w, copier.Apply(line))
			}
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("fieldcopy: read: %w", err)
			}

			fmt.Fprintf(os.Stderr, "fieldcopy: %d line(s) updated\n", copier.Copied())
			return nil
		},
	}

	cmd.Flags().StringVar(&src, "src", "", "source field name to copy from")
	cmd.Flags().StringVar(&dst, "dst", "", "destination field name to copy into")
	return cmd
}
