package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	pctx "github.com/yourorg/logslice/internal/context"
)

func newContextCmd() *cobra.Command {
	var maxLines int
	var since, until string

	cmd := &cobra.Command{
		Use:   "run [file]",
		Short: "Stream lines from a file with max-lines and time-range enforcement",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := pctx.Config{MaxLines: maxLines}

			if since != "" {
				t, err := time.Parse(time.RFC3339, since)
				if err != nil {
					return fmt.Errorf("--since: %w", err)
				}
				cfg.Since = t
			}
			if until != "" {
				t, err := time.Parse(time.RFC3339, until)
				if err != nil {
					return fmt.Errorf("--until: %w", err)
				}
				cfg.Until = t
			}

			pc := pctx.New(cmd.Context(), cfg)

			f, err := os.Open(args[0])
			if err != nil {
				return err
			}
			defer f.Close()

			sc := bufio.NewScanner(f)
			w := bufio.NewWriter(os.Stdout)
			defer w.Flush()

			for sc.Scan() {
				if !pc.RecordLine() {
					break
				}
				fmt.Fprintln(w, sc.Text())
			}

			if !cfg.Quiet {
				fmt.Fprintf(os.Stderr, "lines=%d elapsed=%s\n",
					pc.LinesOut(), pc.Elapsed().Round(time.Millisecond))
			}
			return sc.Err()
		},
	}

	cmd.Flags().IntVar(&maxLines, "max-lines", 0, "stop after N lines (0 = unlimited)")
	cmd.Flags().StringVar(&since, "since", "", "only emit lines after this RFC3339 timestamp")
	cmd.Flags().StringVar(&until, "until", "", "only emit lines before this RFC3339 timestamp")
	cmd.Flags().BoolVar(&pctx.Config{}.Quiet, "quiet", false, "suppress stats output")
	return cmd
}
