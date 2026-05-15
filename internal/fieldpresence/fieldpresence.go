// Package fieldpresence filters log lines based on whether specified fields
// are present or absent in structured (JSON or key=value) log entries.
package fieldpresence

import (
	"fmt"

	"github.com/user/logslice/internal/fieldextract"
)

// Checker filters lines by field presence or absence.
type Checker struct {
	require []string
	forbid  []string
	checked int
	passed  int
}

// New creates a Checker that passes lines containing all fields in require
// and none of the fields in forbid. Either slice may be empty.
func New(require, forbid []string) (*Checker, error) {
	if len(require) == 0 && len(forbid) == 0 {
		return nil, fmt.Errorf("fieldpresence: at least one required or forbidden field must be specified")
	}
	return &Checker{require: require, forbid: forbid}, nil
}

// Allow returns true when the line satisfies all presence constraints.
func (c *Checker) Allow(line string) bool {
	c.checked++
	fields := fieldextract.Extract(line)

	for _, f := range c.require {
		if _, ok := fields[f]; !ok {
			return false
		}
	}
	for _, f := range c.forbid {
		if _, ok := fields[f]; ok {
			return false
		}
	}
	c.passed++
	return true
}

// Checked returns the total number of lines evaluated.
func (c *Checker) Checked() int { return c.checked }

// Passed returns the number of lines that passed the filter.
func (c *Checker) Passed() int { return c.passed }
