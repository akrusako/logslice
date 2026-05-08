// Package redact provides pattern-based line redaction using regular expressions.
// Matching portions of a log line are replaced with a configurable mask string.
package redact

import (
	"fmt"
	"regexp"
)

// Redactor replaces sensitive patterns in log lines.
type Redactor struct {
	patterns []*regexp.Regexp
	mask     string
	count    int64
}

// New creates a Redactor from a slice of regex pattern strings and a mask
// replacement string. Returns an error if any pattern fails to compile.
func New(patterns []string, mask string) (*Redactor, error) {
	if mask == "" {
		mask = "[REDACTED]"
	}
	compiled := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, fmt.Errorf("redact: invalid pattern %q: %w", p, err)
		}
		compiled = append(compiled, re)
	}
	return &Redactor{patterns: compiled, mask: mask}, nil
}

// Apply replaces all matches of any registered pattern in line with the mask.
// Returns the (possibly modified) line and whether any redaction occurred.
func (r *Redactor) Apply(line string) (string, bool) {
	original := line
	for _, re := range r.patterns {
		line = re.ReplaceAllString(line, r.mask)
	}
	if line != original {
		r.count++
		return line, true
	}
	return line, false
}

// Count returns the total number of lines that had at least one redaction.
func (r *Redactor) Count() int64 { return r.count }

// Mask returns the replacement string used by this Redactor.
func (r *Redactor) Mask() string { return r.mask }
