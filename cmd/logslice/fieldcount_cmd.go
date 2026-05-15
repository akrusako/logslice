package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/yourorg/logslice/internal/fieldcount"
)

func newFieldCountCmd() *cobra.Command {
	var min int
	var max int

	cmd := &cobra.Command{
		Use:   "fieldcount <file>",
		Short: "Filter lines by number of structured fields",
		Long: `fieldcount reads a log file and emits only lines whose field count
falls within [--min, --max]. Use --max -1 for no upper bound.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			f, err := fieldcount.New(min, max)
			if err != nil {
				return err
			}

			file, err := os.Open(args[0])
			if err != nil {
				return fmt.Errorf("fieldcount: open %s: %w", args[0], err)
			}
			defer file.Close()

			w := bufio.NewWriter(os.Stdout)
			defer w.Flush()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				if f.Allow(line) {
					fmt.Fprintln(w, line)
				}
			}
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("fieldcount: read: %w", err)
			}

			if verbose, _ := strconv.ParseBool(os.Getenv("LOGSLICE_VERBOSE")); verbose {
				fmt.Fprintf(os.Stderr, "fieldcount: allowed=%d dropped=%d\n",
					f.Allowed(), f.Dropped())
			}
			return nil
		},
	}

	cmd.Flags().IntVar(&min, "min", 0, "minimum number of fields (inclusive)")
	cmd.Flags().IntVar(&max, "max", -1, "maximum number of fields (inclusive), -1 for unlimited")
	return cmd
}
