// Package fieldappend appends a static or derived value to a named field
// in structured log lines (KV or JSON format).
package fieldappend

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Appender appends a suffix to the value of a named field in each log line.
type Appender struct {
	field   string
	suffix  string
	count   int
	missed  int
}

// New creates an Appender that appends suffix to the value of field.
// Returns an error if field or suffix is empty.
func New(field, suffix string) (*Appender, error) {
	if field == "" {
		return nil, fmt.Errorf("fieldappend: field name must not be empty")
	}
	if suffix == "" {
		return nil, fmt.Errorf("fieldappend: suffix must not be empty")
	}
	return &Appender{field: field, suffix: suffix}, nil
}

// Apply appends the configured suffix to the target field value in line.
// Lines where the field is absent are returned unchanged.
func (a *Appender) Apply(line string) string {
	line = strings.TrimRight(line, "\n")
	if line == "" {
		return line
	}

	if strings.HasPrefix(line, "{") {
		result, ok := a.applyJSON(line)
		if ok {
			a.count++
			return result
		}
	} else {
		result, ok := a.applyKV(line)
		if ok {
			a.count++
			return result
		}
	}

	a.missed++
	return line
}

// Count returns the number of lines where the field was found and modified.
func (a *Appender) Count() int { return a.count }

// Missed returns the number of lines where the target field was not found.
func (a *Appender) Missed() int { return a.missed }

func (a *Appender) applyJSON(line string) (string, bool) {
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		return line, false
	}
	v, ok := m[a.field]
	if !ok {
		return line, false
	}
	m[a.field] = fmt.Sprintf("%v%s", v, a.suffix)
	b, err := json.Marshal(m)
	if err != nil {
		return line, false
	}
	return string(b), true
}

func (a *Appender) applyKV(line string) (string, bool) {
	parts := strings.Fields(line)
	modified := false
	for i, p := range parts {
		if idx := strings.IndexByte(p, '='); idx > 0 {
			key := p[:idx]
			val := p[idx+1:]
			if key == a.field {
				parts[i] = key + "=" + val + a.suffix
				modified = true
			}
		}
	}
	if !modified {
		return line, false
	}
	return strings.Join(parts, " "), true
}
