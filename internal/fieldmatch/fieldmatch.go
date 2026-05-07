// Package fieldmatch provides regex-based matching against structured log fields.
package fieldmatch

import (
	"fmt"
	"regexp"

	"github.com/user/logslice/internal/fieldextract"
)

// Matcher holds a compiled set of field-pattern rules.
type Matcher struct {
	rules []rule
}

type rule struct {
	field string
	re    *regexp.Regexp
}

// New creates a Matcher from a map of field name to regex pattern strings.
// Returns an error if any pattern fails to compile.
func New(patterns map[string]string) (*Matcher, error) {
	m := &Matcher{}
	for field, pat := range patterns {
		re, err := regexp.Compile(pat)
		if err != nil {
			return nil, fmt.Errorf("fieldmatch: field %q pattern %q: %w", field, pat, err)
		}
		m.rules = append(m.rules, rule{field: field, re: re})
	}
	return m, nil
}

// Match returns true when every rule's field value in the line matches its
// compiled regex. Lines that do not contain a required field do not match.
// If the Matcher has no rules it matches every line.
func (m *Matcher) Match(line string) bool {
	if len(m.rules) == 0 {
		return true
	}
	fields := fieldextract.Extract(line)
	for _, r := range m.rules {
		v, ok := fields[r.field]
		if !ok {
			return false
		}
		if !r.re.MatchString(v) {
			return false
		}
	}
	return true
}

// MatchedFields returns only the fields that satisfy at least one rule.
// Useful for reporting which fields triggered a match.
func (m *Matcher) MatchedFields(line string) map[string]string {
	result := make(map[string]string)
	fields := fieldextract.Extract(line)
	for _, r := range m.rules {
		v, ok := fields[r.field]
		if ok && r.re.MatchString(v) {
			result[r.field] = v
		}
	}
	return result
}

// Fields returns the list of field names that this Matcher has rules for.
// The order of the returned slice is not guaranteed to be stable.
func (m *Matcher) Fields() []string {
	fields := make([]string, len(m.rules))
	for i, r := range m.rules {
		fields[i] = r.field
	}
	return fields
}
