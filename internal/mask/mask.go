// Package mask provides redaction of sensitive fields in log lines.
// It supports both JSON-structured logs and key=value formatted lines,
// replacing matched field values with a configurable placeholder.
package mask

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

const defaultPlaceholder = "***"

// Masker redacts specified field values from log lines.
type Masker struct {
	fields      map[string]struct{}
	placeholder string
	kvPattern   *regexp.Regexp
	masked      int64
}

// New creates a Masker that will redact the given field names.
// If placeholder is empty, "***" is used.
func New(fields []string, placeholder string) *Masker {
	if placeholder == "" {
		placeholder = defaultPlaceholder
	}
	set := make(map[string]struct{}, len(fields))
	for _, f := range fields {
		set[strings.TrimSpace(f)] = struct{}{}
	}
	return &Masker{
		fields:      set,
		placeholder: placeholder,
		kvPattern:   regexp.MustCompile(`(\w+)=([^\s]+)`),
	}
}

// Mask redacts sensitive fields from line and returns the result.
func (m *Masker) Mask(line string) string {
	if len(m.fields) == 0 {
		return line
	}
	trimmed := strings.TrimSpace(line)
	if strings.HasPrefix(trimmed, "{") {
		if out, ok := m.maskJSON(trimmed); ok {
			m.masked++
			return out
		}
	}
	out, changed := m.maskKV(line)
	if changed {
		m.masked++
	}
	return out
}

// MaskedCount returns the total number of lines that had at least one field redacted.
func (m *Masker) MaskedCount() int64 {
	return m.masked
}

func (m *Masker) maskJSON(line string) (string, bool) {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line, false
	}
	changed := false
	for k := range m.fields {
		if _, ok := obj[k]; ok {
			obj[k] = m.placeholder
			changed = true
		}
	}
	if !changed {
		return line, false
	}
	b, err := json.Marshal(obj)
	if err != nil {
		return line, false
	}
	return string(b), true
}

func (m *Masker) maskKV(line string) (string, bool) {
	changed := false
	out := m.kvPattern.ReplaceAllStringFunc(line, func(match string) string {
		parts := strings.SplitN(match, "=", 2)
		if len(parts) != 2 {
			return match
		}
		if _, ok := m.fields[parts[0]]; ok {
			changed = true
			return fmt.Sprintf("%s=%s", parts[0], m.placeholder)
		}
		return match
	})
	return out, changed
}
