package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/yourorg/logslice/internal/fieldclip"
	"github.com/spf13/cobra"
)

func newFieldClipCmd() *cobra.Command {
	var minVal, maxVal float64
	var field string

	cmd := &cobra.Command{
		Use:   "fieldclip <file>",
		Short: "Clamp a numeric field to [min, max] in every log line",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if field == "" {
				return fmt.Errorf("--field is required")
			}

			clipper, err := fieldclip.New(field, minVal, maxVal)
			if err != nil {
				return err
			}

			f, err := os.Open(args[0])
			if err != nil {
				return fmt.Errorf("fieldclip: open %s: %w", args[0], err)
			}
			defer f.Close()

			w := bufio.NewWriter(os.Stdout)
			defer w.Flush()

			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				line := scanner.Text()
				fmt.Fprintln(w, clipper.Apply(line))
			}
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("fieldclip: read: %w", err)
			}

			fmt.Fprintf(os.Stderr, "fieldclip: %d value(s) clamped\n",
				clipper.Clipped())
			return nil
		},
	}

	cmd.Flags().StringVar(&field, "field", "", "field name to clamp (required)")
	cmd.Flags().Float64Var(&minVal, "min", 0,
		"minimum allowed value (default 0); use "+strconv.FormatFloat(-1e308, 'g', -1, 64)+" for no lower bound")
	cmd.Flags().Float64Var(&maxVal, "max", 0,
		"maximum allowed value (required)")
	_ = cmd.MarkFlagRequired("max")

	return cmd
}
