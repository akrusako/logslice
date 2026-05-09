// Package prefix prepends a static or dynamic string to each log line.
package prefix

import "strings"

// Prefixer prepends a label to every line it processes.
type Prefixer struct {
	prefix  string
	count   int
	dropped int
}

// New returns a Prefixer that prepends p to each non-empty line.
// If p is empty the Prefixer is effectively a no-op.
func New(p string) *Prefixer {
	return &Prefixer{prefix: p}
}

// Apply prepends the configured prefix to line and returns the result.
// Empty lines are passed through unchanged and counted as dropped.
func (p *Prefixer) Apply(line string) string {
	if line == "" {
		p.dropped++
		return line
	}
	p.count++
	if p.prefix == "" {
		return line
	}
	return p.prefix + line
}

// Count returns the number of non-empty lines processed.
func (p *Prefixer) Count() int { return p.count }

// Dropped returns the number of empty lines skipped.
func (p *Prefixer) Dropped() int { return p.dropped }

// SetPrefix replaces the current prefix at runtime.
func (p *Prefixer) SetPrefix(s string) { p.prefix = s }

// Strip removes the prefix from line if it is present and returns the
// stripped string along with a boolean indicating whether the prefix was found.
func (p *Prefixer) Strip(line string) (string, bool) {
	if p.prefix == "" {
		return line, true
	}
	if strings.HasPrefix(line, p.prefix) {
		return strings.TrimPrefix(line, p.prefix), true
	}
	return line, false
}
