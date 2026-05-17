// Package fieldsplit splits a single field value into multiple fields
// using a configurable delimiter.
package fieldsplit

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Splitter splits a named field's value into indexed sub-fields.
type Splitter struct {
	field     string
	delim     string
	prefix    string
	split     int64
	dropped   int64
}

// New creates a Splitter that splits field by delim, naming the resulting
// sub-fields <prefix>0, <prefix>1, …  A non-positive maxParts means no limit.
func New(field, delim, prefix string) (*Splitter, error) {
	if field == "" {
		return nil, fmt.Errorf("fieldsplit: field must not be empty")
	}
	if delim == "" {
		return nil, fmt.Errorf("fieldsplit: delim must not be empty")
	}
	if prefix == "" {
		prefix = field + "_"
	}
	return &Splitter{field: field, delim: delim, prefix: prefix}, nil
}

// Apply returns the line with the target field split into indexed sub-fields.
// If the field is absent or the line is plain text the original line is
// returned unchanged.
func (s *Splitter) Apply(line string) string {
	if strings.HasPrefix(strings.TrimSpace(line), "{") {
		return s.applyJSON(line)
	}
	return s.applyKV(line)
}

func (s *Splitter) applyJSON(line string) string {
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		s.dropped++
		return line
	}
	val, ok := m[s.field]
	if !ok {
		return line
	}
	str, ok := val.(string)
	if !ok {
		return line
	}
	parts := strings.Split(str, s.delim)
	delete(m, s.field)
	for i, p := range parts {
		m[fmt.Sprintf("%s%d", s.prefix, i)] = p
	}
	b, err := json.Marshal(m)
	if err != nil {
		s.dropped++
		return line
	}
	s.split++
	return string(b)
}

func (s *Splitter) applyKV(line string) string {
	pairs := strings.Fields(line)
	var out []string
	matched := false
	for _, pair := range pairs {
		idx := strings.IndexByte(pair, '=')
		if idx < 0 {
			out = append(out, pair)
			continue
		}
		key, val := pair[:idx], pair[idx+1:]
		if key != s.field {
			out = append(out, pair)
			continue
		}
		parts := strings.Split(val, s.delim)
		for i, p := range parts {
			out = append(out, fmt.Sprintf("%s%d=%s", s.prefix, i, p))
		}
		matched = true
	}
	if !matched {
		return line
	}
	s.split++
	return strings.Join(out, " ")
}

// Split returns the number of lines that had the field split.
func (s *Splitter) Split() int64 { return s.split }

// Dropped returns the number of lines that could not be processed.
func (s *Splitter) Dropped() int64 { return s.dropped }
