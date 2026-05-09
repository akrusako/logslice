package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/yourorg/logslice/internal/prefix"
	"github.com/spf13/cobra"
)

func newPrefixCmd() *cobra.Command {
	var label string
	var strip bool

	cmd := &cobra.Command{
		Use:   "prefix <file>",
		Short: "Prepend (or strip) a label from every log line",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			f, err := os.Open(args[0])
			if err != nil {
				return fmt.Errorf("prefix: open file: %w", err)
			}
			defer f.Close()

			p := prefix.New(label)
			scanner := bufio.NewScanner(f)
			w := bufio.NewWriter(os.Stdout)
			defer w.Flush()

			for scanner.Scan() {
				line := scanner.Text()
				var out string
				if strip {
					stripped, _ := p.Strip(line)
					out = stripped
				} else {
					out = p.Apply(line)
				}
				fmt.Fprintln(w, out)
			}
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("prefix: scan: %w", err)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&label, "label", "l", "", "prefix string to prepend to each line")
	cmd.Flags().BoolVarP(&strip, "strip", "s", false, "strip the prefix instead of adding it")
	return cmd
}
