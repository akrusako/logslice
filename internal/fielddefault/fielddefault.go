// Package fielddefault fills in missing log fields with a configured default value.
package fielddefault

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Defaulter fills missing fields with a default value.
type Defaulter struct {
	field        string
	defaultValue string
	applied      int64
	skipped      int64
}

// New creates a Defaulter that sets field to defaultValue when the field is absent.
func New(field, defaultValue string) (*Defaulter, error) {
	if field == "" {
		return nil, fmt.Errorf("fielddefault: field name must not be empty")
	}
	return &Defaulter{field: field, defaultValue: defaultValue}, nil
}

// Apply returns the line with the field injected if it was missing, or the
// original line if the field already exists.
func (d *Defaulter) Apply(line string) string {
	line = strings.TrimRight(line, "\n")
	if line == "" {
		return line
	}
	if strings.HasPrefix(line, "{") {
		result := d.applyJSON(line)
		return result
	}
	return d.applyKV(line)
}

func (d *Defaulter) applyJSON(line string) string {
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		return line
	}
	if _, ok := m[d.field]; ok {
		d.skipped++
		return line
	}
	m[d.field] = d.defaultValue
	b, err := json.Marshal(m)
	if err != nil {
		return line
	}
	d.applied++
	return string(b)
}

func (d *Defaulter) applyKV(line string) string {
	prefix := d.field + "="
	for _, part := range strings.Fields(line) {
		if strings.HasPrefix(part, prefix) {
			d.skipped++
			return line
		}
	}
	d.applied++
	return line + " " + d.field + "=" + d.defaultValue
}

// Applied returns the number of lines where the default was injected.
func (d *Defaulter) Applied() int64 { return d.applied }

// Skipped returns the number of lines that already had the field.
func (d *Defaulter) Skipped() int64 { return d.skipped }
