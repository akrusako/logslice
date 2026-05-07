// Package grep provides regex-based line matching for log pipelines.
package grep

import (
	"fmt"
	"regexp"
)

// Matcher holds a compiled set of patterns and applies them to log lines.
type Matcher struct {
	patterns []*regexp.Regexp
	invert   bool
	matched  int64
	dropped  int64
}

// New compiles the given regex patterns and returns a Matcher.
// If invert is true, lines that match any pattern are dropped instead of kept.
func New(patterns []string, invert bool) (*Matcher, error) {
	if len(patterns) == 0 {
		return &Matcher{invert: invert}, nil
	}
	compiled := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, fmt.Errorf("grep: invalid pattern %q: %w", p, err)
		}
		compiled = append(compiled, re)
	}
	return &Matcher{patterns: compiled, invert: invert}, nil
}

// Match reports whether the line should be kept.
func (m *Matcher) Match(line string) bool {
	if len(m.patterns) == 0 {
		m.matched++
		return true
	}
	hit := false
	for _, re := range m.patterns {
		if re.MatchString(line) {
			hit = true
			break
		}
	}
	keep := hit != m.invert // XOR: invert flips the decision
	if keep {
		m.matched++
	} else {
		m.dropped++
	}
	return keep
}

// Matched returns the number of lines kept.
func (m *Matcher) Matched() int64 { return m.matched }

// Dropped returns the number of lines rejected.
func (m *Matcher) Dropped() int64 { return m.dropped }
