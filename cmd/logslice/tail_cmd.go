package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/output"
	"github.com/yourorg/logslice/internal/tail"
)

// tailCmd implements the `logslice tail` sub-command.
type tailCmd struct {
	n      int
	fmt    string
	substr string
	fields []string
}

func (c *tailCmd) register(fs *flag.FlagSet) {
	fs.IntVar(&c.n, "n", 10, "number of lines to show from the end")
	fs.StringVar(&c.fmt, "format", "raw", "output format: raw|json|tsv")
	fs.StringVar(&c.substr, "match", "", "substring filter")
}

func (c *tailCmd) run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("tail: no file specified")
	}
	path := args[0]

	r, err := tail.New(path)
	if err != nil {
		return err
	}

	lines, err := r.Lines(c.n)
	if err != nil {
		return err
	}

	f, err := filter.NewFromArgs(c.substr, c.fields)
	if err != nil {
		return err
	}

	out := output.New(os.Stdout, c.fmt)

	for _, l := range lines {
		if !f.Match(l) {
			continue
		}
		if err := out.Write(l); err != nil {
			return err
		}
	}
	return out.Flush()
}
