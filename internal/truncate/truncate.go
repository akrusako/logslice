// Package truncate provides line truncation for log output,
// trimming lines that exceed a configurable maximum byte length.
package truncate

import "fmt"

const defaultSuffix = "..."

// Truncator trims log lines that exceed a maximum length.
type Truncator struct {
	maxLen int
	suffix string
	truncated int
}

// New returns a Truncator that trims lines longer than maxLen bytes.
// If maxLen is zero or negative, no truncation is applied.
func New(maxLen int) *Truncator {
	return &Truncator{
		maxLen: maxLen,
		suffix: defaultSuffix,
	}
}

// Apply returns the (possibly truncated) form of line.
// If maxLen <= 0, the original line is returned unchanged.
func (t *Truncator) Apply(line string) string {
	if t.maxLen <= 0 || len(line) <= t.maxLen {
		return line
	}
	t.truncated++
	cutAt := t.maxLen - len(t.suffix)
	if cutAt < 0 {
		cutAt = 0
	}
	return line[:cutAt] + t.suffix
}

// TruncatedCount returns the total number of lines that have been truncated.
func (t *Truncator) TruncatedCount() int {
	return t.truncated
}

// Summary returns a human-readable summary of truncation activity.
func (t *Truncator) Summary() string {
	if t.truncated == 0 {
		return "no lines truncated"
	}
	return fmt.Sprintf("%d line(s) truncated to %d bytes", t.truncated, t.maxLen)
}
