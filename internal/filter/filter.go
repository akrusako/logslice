// Package filter provides structured log line filtering based on
// key=value field matching and substring search.
package filter

import (
	"strings"
)

// Filter holds compiled filter criteria for matching log lines.
type Filter struct {
	fields    map[string]string
	substring string
}

// Options configures a Filter.
type Options struct {
	// Fields is a map of key=value pairs that must all be present in a log line.
	Fields map[string]string
	// Substring requires this string to appear anywhere in the log line.
	Substring string
}

// New creates a new Filter from the given Options.
func New(opts Options) *Filter {
	fields := make(map[string]string, len(opts.Fields))
	for k, v := range opts.Fields {
		fields[k] = v
	}
	return &Filter{
		fields:    fields,
		substring: opts.Substring,
	}
}

// Match reports whether line satisfies all filter criteria.
// An empty filter matches every line.
func (f *Filter) Match(line string) bool {
	if f.substring != "" && !strings.Contains(line, f.substring) {
		return false
	}
	for k, v := range f.fields {
		pair := k + "=" + v
		if !strings.Contains(line, pair) {
			return false
		}
	}
	return true
}

// Empty reports whether the filter has no criteria (matches everything).
func (f *Filter) Empty() bool {
	return f.substring == "" && len(f.fields) == 0
}
