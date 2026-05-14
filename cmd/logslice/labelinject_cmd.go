package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/user/logslice/internal/labelinject"
	"github.com/spf13/cobra"
)

func newLabelInjectCmd() *cobra.Command {
	var labels []string

	cmd := &cobra.Command{
		Use:   "labelinject <file>",
		Short: "Inject static labels into every log line",
		Long: `Read a log file and append static key=value labels to each line.

Labels are injected in a format-aware way: JSON lines receive new top-level
fields; key=value lines and plain text receive appended key=value pairs.

Example:
  logslice labelinject --label env=prod --label region=us-east app.log`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(labels) == 0 {
				return fmt.Errorf("at least one --label is required")
			}

			inj, err := labelinject.New(labels)
			if err != nil {
				return err
			}

			f, err := os.Open(args[0])
			if err != nil {
				return fmt.Errorf("open %s: %w", args[0], err)
			}
			defer f.Close()

			scanner := bufio.NewScanner(f)
			writer := bufio.NewWriter(os.Stdout)
			defer writer.Flush()

			for scanner.Scan() {
				line := scanner.Text()
				fmt.Fprintln(writer, inj.Apply(line))
			}
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("scan: %w", err)
			}

			fmt.Fprintf(os.Stderr, "labelinject: %d lines labelled\n", inj.Injected())
			return nil
		},
	}

	cmd.Flags().StringArrayVar(&labels, "label", nil, "label to inject as key=value (repeatable)")
	return cmd
}
