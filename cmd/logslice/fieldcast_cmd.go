package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourorg/logslice/internal/fieldcast"
)

func newFieldCastCmd() *cobra.Command {
	var castType string

	cmd := &cobra.Command{
		Use:   "fieldcast <file> <field>",
		Short: "Coerce a log field to a target type (string|int|float|bool)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]
			field := args[1]

			to, err := fieldcast.ParseType(castType)
			if err != nil {
				return err
			}

			caster, err := fieldcast.New(field, to)
			if err != nil {
				return err
			}

			f, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("fieldcast: open %s: %w", path, err)
			}
			defer f.Close()

			scanner := bufio.NewScanner(f)
			writer := bufio.NewWriter(os.Stdout)
			defer writer.Flush()

			for scanner.Scan() {
				line := scanner.Text()
				fmt.Fprintln(writer, caster.Apply(line))
			}
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("fieldcast: read %s: %w", path, err)
			}

			fmt.Fprintf(os.Stderr, "fieldcast: casted=%d failed=%d\n",
				caster.Casted(), caster.Failed())
			return nil
		},
	}

	cmd.Flags().StringVar(&castType, "type", "string",
		"Target type: string, int, float, bool")
	return cmd
}
