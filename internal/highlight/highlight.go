// Package highlight provides terminal colour highlighting for matched
// substrings and field values within log lines.
package highlight

import (
	"strings"
)

// ANSI escape codes.
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Cyan   = "\033[36m"
	Bold   = "\033[1m"
)

// Highlighter wraps matched terms in ANSI colour codes.
type Highlighter struct {
	colour  string
	enabled bool
}

// New returns a Highlighter. When enabled is false all methods return the
// input unchanged, making it safe to use unconditionally in pipelines.
func New(colour string, enabled bool) *Highlighter {
	if colour == "" {
		colour = Yellow
	}
	return &Highlighter{colour: colour, enabled: enabled}
}

// Term wraps every occurrence of term in line with the configured colour.
// The comparison is case-sensitive.
func (h *Highlighter) Term(line, term string) string {
	if !h.enabled || term == "" {
		return line
	}
	return strings.ReplaceAll(line, term, h.colour+Bold+term+Reset)
}

// Terms applies Term for each entry in terms sequentially.
func (h *Highlighter) Terms(line string, terms []string) string {
	if !h.enabled {
		return line
	}
	for _, t := range terms {
		line = h.Term(line, t)
	}
	return line
}

// Field highlights the value portion of a key=value or "key":"value" pair.
// It wraps the raw value string wherever it appears in line.
func (h *Highlighter) Field(line, value string) string {
	return h.Term(line, value)
}
