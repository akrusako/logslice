package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/yourorg/logslice/internal/fieldpad"
	"github.com/spf13/cobra"
)

func newFieldPadCmd() *cobra.Command {
	var (
		field    string
		width    int
		padChar  string
		dirFlag  string
	)

	cmd := &cobra.Command{
		Use:   "fieldpad <file>",
		Short: "Pad a log field value to a fixed column width",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if field == "" {
				return fmt.Errorf("--field is required")
			}

			var ch rune = ' '
			if padChar != "" {
				runes := []rune(padChar)
				if len(runes) != 1 {
					return fmt.Errorf("--char must be a single character")
				}
				ch = runes[0]
			}

			dir := fieldpad.Right
			if strings.EqualFold(dirFlag, "left") {
				dir = fieldpad.Left
			}

			padder, err := fieldpad.New(field, width, ch, dir)
			if err != nil {
				return err
			}

			f, err := os.Open(args[0])
			if err != nil {
				return err
			}
			defer f.Close()

			scanner := bufio.NewScanner(f)
			w := bufio.NewWriter(os.Stdout)
			defer w.Flush()

			for scanner.Scan() {
				fmt.Fprintln(w, padder.Apply(scanner.Text()))
			}
			if err := scanner.Err(); err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "padded %d lines\n", padder.PaddedCount())
			return nil
		},
	}

	cmd.Flags().StringVar(&field, "field", "", "field name to pad (required)")
	cmd.Flags().IntVar(&width, "width", 10, "target column width")
	cmd.Flags().StringVar(&padChar, "char", " ", "padding character (single rune)")
	cmd.Flags().StringVar(&dirFlag, "dir", "right", "padding direction: left or right")
	return cmd
}
