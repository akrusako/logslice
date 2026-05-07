// Package fieldrewrite provides a processor that rewrites field values
// in structured log lines using simple find-and-replace rules.
package fieldrewrite

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Rule describes a single field rewrite: for the given field name, replace
// occurrences of Old with New in the field's string value.
type Rule struct {
	Field string
	Old   string
	New   string
}

// Rewriter applies a set of rewrite rules to log lines.
type Rewriter struct {
	rules []Rule
}

// New creates a Rewriter with the provided rules. An error is returned if any
// rule has an empty Field name.
func New(rules []Rule) (*Rewriter, error) {
	for i, r := range rules {
		if r.Field == "" {
			return nil, fmt.Errorf("rule %d: field name must not be empty", i)
		}
	}
	return &Rewriter{rules: rules}, nil
}

// Apply rewrites the given log line according to the configured rules.
// JSON and key=value formats are both supported. Lines that cannot be
// parsed are returned unchanged.
func (rw *Rewriter) Apply(line string) string {
	if len(rw.rules) == 0 {
		return line
	}
	if rewritten, ok := rw.applyJSON(line); ok {
		return rewritten
	}
	if rewritten, ok := rw.applyKV(line); ok {
		return rewritten
	}
	return line
}

func (rw *Rewriter) applyJSON(line string) (string, bool) {
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		return "", false
	}
	modified := false
	for _, r := range rw.rules {
		if v, ok := m[r.Field]; ok {
			if s, ok := v.(string); ok {
				newVal := strings.ReplaceAll(s, r.Old, r.New)
				if newVal != s {
					m[r.Field] = newVal
					modified = true
				}
			}
		}
	}
	if !modified {
		return line, true
	}
	b, err := json.Marshal(m)
	if err != nil {
		return line, true
	}
	return string(b), true
}

func (rw *Rewriter) applyKV(line string) (string, bool) {
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return line, false
	}
	hasKV := false
	for _, p := range parts {
		if strings.Contains(p, "=") {
			hasKV = true
			break
		}
	}
	if !hasKV {
		return "", false
	}
	for i, p := range parts {
		idx := strings.IndexByte(p, '=')
		if idx < 0 {
			continue
		}
		key := p[:idx]
		val := p[idx+1:]
		for _, r := range rw.rules {
			if key == r.Field {
				val = strings.ReplaceAll(val, r.Old, r.New)
			}
		}
		parts[i] = key + "=" + val
	}
	return strings.Join(parts, " "), true
}
