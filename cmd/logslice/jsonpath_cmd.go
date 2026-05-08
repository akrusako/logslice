package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/user/logslice/internal/jsonpath"
	"github.com/spf13/cobra"
)

// newJSONPathCmd returns a cobra command that extracts a dot-notation field
// value from every JSON log line in a file, printing matched lines.
func newJSONPathCmd() *cobra.Command {
	var (
		pathFlag  string
		valueFlag string
		invert    bool
	)

	cmd := &cobra.Command{
		Use:   "jsonpath <file>",
		Short: "Filter log lines by a dot-notation JSON field value",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if pathFlag == "" {
				return fmt.Errorf("--path is required")
			}
			ext, err := jsonpath.New(pathFlag)
			if err != nil {
				return fmt.Errorf("invalid path: %w", err)
			}
			f, err := os.Open(args[0])
			if err != nil {
				return fmt.Errorf("open: %w", err)
			}
			defer f.Close()

			scanner := bufio.NewScanner(f)
			w := bufio.NewWriter(os.Stdout)
			defer w.Flush()

			for scanner.Scan() {
				line := scanner.Text()
				val, ok := ext.Get(line)
				matched := ok && (valueFlag == "" || val == valueFlag)
				if invert {
					matched = !matched
				}
				if matched {
					fmt.Fprintln(w, line)
				}
			}
			return scanner.Err()
		},
	}

	cmd.Flags().StringVar(&pathFlag, "path", "", "dot-notation field path (e.g. request.method)")
	cmd.Flags().StringVar(&valueFlag, "value", "", "required field value (omit to match any)")
	cmd.Flags().BoolVar(&invert, "invert", false, "invert match (exclude matching lines)")
	return cmd
}
