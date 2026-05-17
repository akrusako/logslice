// Package fieldformat applies a sprintf-style format template to a named
// log field, replacing its value with the formatted result.
package fieldformat

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Formatter rewrites a single field using a Go format string.
type Formatter struct {
	field   string
	fmt     string
	count   int64
}

// New returns a Formatter that rewrites field using the given Go format
// string (e.g. "%.2f" or "%08d"). Returns an error if field or format
// is empty.
func New(field, format string) (*Formatter, error) {
	if field == "" {
		return nil, fmt.Errorf("fieldformat: field must not be empty")
	}
	if format == "" {
		return nil, fmt.Errorf("fieldformat: format must not be empty")
	}
	return &Formatter{field: field, fmt: format}, nil
}

// Apply rewrites the target field in line and returns the result.
// Lines where the field is absent are returned unchanged.
func (f *Formatter) Apply(line string) string {
	line = strings.TrimRight(line, "\n")
	if line == "" {
		return line
	}
	if strings.HasPrefix(strings.TrimSpace(line), "{") {
		result := f.applyJSON(line)
		if result != line {
			f.count++
		}
		return result
	}
	result := f.applyKV(line)
	if result != line {
		f.count++
	}
	return result
}

func (f *Formatter) applyJSON(line string) string {
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		return line
	}
	val, ok := m[f.field]
	if !ok {
		return line
	}
	m[f.field] = fmt.Sprintf(f.fmt, val)
	b, err := json.Marshal(m)
	if err != nil {
		return line
	}
	return string(b)
}

func (f *Formatter) applyKV(line string) string {
	parts := strings.Fields(line)
	for i, p := range parts {
		if !strings.Contains(p, "=") {
			continue
		}
		idx := strings.IndexByte(p, '=')
		key := p[:idx]
		val := p[idx+1:]
		if key != f.field {
			continue
		}
		parts[i] = key + "=" + fmt.Sprintf(f.fmt, val)
		return strings.Join(parts, " ")
	}
	return line
}

// FormattedCount returns the number of lines where the field was rewritten.
func (f *Formatter) FormattedCount() int64 { return f.count }
