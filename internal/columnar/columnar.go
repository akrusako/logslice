// Package columnar extracts and formats specific fields from log lines
// as aligned columns, suitable for tabular display or downstream processing.
package columnar

import (
	"fmt"
	"strings"

	"github.com/user/logslice/internal/fieldextract"
)

// Columnar holds configuration for column extraction.
type Columnar struct {
	fields  []string
	width    int
	sep      string
	dropped  int
}

// New creates a Columnar extractor for the given field names.
// width is the minimum column width for padding (0 disables padding).
// sep is the column separator string.
func New(fields []string, width int, sep string) *Columnar {
	f := make([]string, len(fields))
	copy(f, fields)
	if sep == "" {
		sep = "\t"
	}
	return &Columnar{fields: f, width: width, sep: sep}
}

// Format extracts the configured fields from line and returns a
// column-formatted string. If a field is missing the column is left blank.
// Returns the formatted row and whether all fields were found.
func (c *Columnar) Format(line string) (string, bool) {
	if len(c.fields) == 0 {
		return line, true
	}

	kv := fieldextract.Extract(line)
	cols := make([]string, len(c.fields))
	allFound := true

	for i, f := range c.fields {
		v, ok := kv[f]
		if !ok {
			allFound = false
			v = ""
		}
		if c.width > 0 {
			cols[i] = fmt.Sprintf("%-*s", c.width, v)
		} else {
			cols[i] = v
		}
	}

	if !allFound {
		c.dropped++
	}

	return strings.Join(cols, c.sep), allFound
}

// Header returns a formatted header row using the field names as column labels.
func (c *Columnar) Header() string {
	cols := make([]string, len(c.fields))
	for i, f := range c.fields {
		if c.width > 0 {
			cols[i] = fmt.Sprintf("%-*s", c.width, f)
		} else {
			cols[i] = f
		}
	}
	return strings.Join(cols, c.sep)
}

// DroppedCount returns the number of lines where at least one field was absent.
func (c *Columnar) DroppedCount() int {
	return c.dropped
}
