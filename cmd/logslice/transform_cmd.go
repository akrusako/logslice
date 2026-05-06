package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/yourorg/logslice/internal/transform"
	"github.com/spf13/cobra"
)

func newTransformCmd() *cobra.Command {
	var (
		upper       bool
		lower       bool
		stripPrefix string
		stripSuffix string
		renameRaw   []string
	)

	cmd := &cobra.Command{
		Use:   "transform <file>",
		Short: "Apply text transformations to each log line",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			f, err := os.Open(args[0])
			if err != nil {
				return fmt.Errorf("open: %w", err)
			}
			defer f.Close()

			var opts []transform.Option
			if upper {
				opts = append(opts, transform.WithUpper())
			}
			if lower {
				opts = append(opts, transform.WithLower())
			}
			if stripPrefix != "" {
				opts = append(opts, transform.WithStripPrefix(stripPrefix))
			}
			if stripSuffix != "" {
				opts = append(opts, transform.WithStripSuffix(stripSuffix))
			}
			if len(renameRaw) > 0 {
				m := make(map[string]string, len(renameRaw))
				for _, r := range renameRaw {
					parts := strings.SplitN(r, "=", 2)
					if len(parts) != 2 {
						return fmt.Errorf("invalid rename %q: want old=new", r)
					}
					m[parts[0]] = parts[1]
				}
				opts = append(opts, transform.WithRenameFields(m))
			}

			tr := transform.New(opts...)
			scanner := bufio.NewScanner(f)
			w := bufio.NewWriter(os.Stdout)
			defer w.Flush()
			for scanner.Scan() {
				fmt.Fprintln(w, tr.Apply(scanner.Text()))
			}
			return scanner.Err()
		},
	}

	cmd.Flags().BoolVar(&upper, "upper", false, "Convert lines to upper case")
	cmd.Flags().BoolVar(&lower, "lower", false, "Convert lines to lower case")
	cmd.Flags().StringVar(&stripPrefix, "strip-prefix", "", "Strip prefix from each line")
	cmd.Flags().StringVar(&stripSuffix, "strip-suffix", "", "Strip suffix from each line")
	cmd.Flags().StringArrayVar(&renameRaw, "rename", nil, "Rename field old=new (repeatable)")
	return cmd
}
