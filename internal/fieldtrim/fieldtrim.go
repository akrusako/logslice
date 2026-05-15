// Package fieldtrim trims whitespace (or a custom cutset) from the value
// of a named field in structured log lines.
package fieldtrim

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Trimmer trims a field value in structured log lines.
type Trimmer struct {
	field   string
	cutset  string
	trimmed int
}

// New returns a Trimmer that trims the given field.
// cutset is the set of characters to remove from both ends of the value;
// pass an empty string to trim Unicode whitespace.
func New(field, cutset string) (*Trimmer, error) {
	if field == "" {
		return nil, fmt.Errorf("fieldtrim: field name must not be empty")
	}
	return &Trimmer{field: field, cutset: cutset}, nil
}

// Apply trims the field value in line and returns the result.
func (t *Trimmer) Apply(line string) string {
	if line == "" {
		return line
	}
	var out string
	if strings.HasPrefix(strings.TrimSpace(line), "{") {
		out = t.applyJSON(line)
	} else {
		out = t.applyKV(line)
	}
	if out != line {
		t.trimmed++
	}
	return out
}

// Trimmed returns the number of lines where a field value was changed.
func (t *Trimmer) Trimmed() int { return t.trimmed }

func (t *Trimmer) trim(s string) string {
	if t.cutset == "" {
		return strings.TrimSpace(s)
	}
	return strings.Trim(s, t.cutset)
}

func (t *Trimmer) applyJSON(line string) string {
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		return line
	}
	v, ok := m[t.field]
	if !ok {
		return line
	}
	s, ok := v.(string)
	if !ok {
		return line
	}
	m[t.field] = t.trim(s)
	b, err := json.Marshal(m)
	if err != nil {
		return line
	}
	return string(b)
}

func (t *Trimmer) applyKV(line string) string {
	parts := strings.Fields(line)
	changed := false
	for i, p := range parts {
		if !strings.Contains(p, "=") {
			continue
		}
		idx := strings.IndexByte(p, '=')
		key := p[:idx]
		val := p[idx+1:]
		if key != t.field {
			continue
		}
		trimmed := t.trim(val)
		if trimmed != val {
			parts[i] = key + "=" + trimmed
			changed = true
		}
	}
	if !changed {
		return line
	}
	return strings.Join(parts, " ")
}
