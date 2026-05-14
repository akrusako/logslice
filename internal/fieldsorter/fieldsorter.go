// Package fieldsorter reorders fields in structured log lines.
// It supports JSON and key=value formats, placing named fields
// first in the specified order, with remaining fields appended.
package fieldsorter

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Sorter reorders fields in structured log lines.
type Sorter struct {
	order   []string
	sorted  uint64
	skipped uint64
}

// New creates a Sorter that places fields in the given order.
// Fields not listed appear after the ordered fields.
func New(order []string) (*Sorter, error) {
	if len(order) == 0 {
		return nil, fmt.Errorf("fieldsorter: at least one field required")
	}
	return &Sorter{order: order}, nil
}

// Apply reorders fields in line and returns the result.
func (s *Sorter) Apply(line string) string {
	line = strings.TrimSpace(line)
	if line == "" {
		s.skipped++
		return line
	}
	if strings.HasPrefix(line, "{") {
		if out, ok := s.sortJSON(line); ok {
			s.sorted++
			return out
		}
	}
	if out, ok := s.sortKV(line); ok {
		s.sorted++
		return out
	}
	s.skipped++
	return line
}

// Sorted returns the number of lines whose fields were reordered.
func (s *Sorter) Sorted() uint64 { return s.sorted }

// Skipped returns the number of lines that were not reordered.
func (s *Sorter) Skipped() uint64 { return s.skipped }

func (s *Sorter) sortJSON(line string) (string, bool) {
	var m map[string]json.RawMessage
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		return "", false
	}
	out := make(map[string]json.RawMessage, len(m))
	for k, v := range m {
		out[k] = v
	}
	ordered := make([]string, 0, len(m))
	seen := map[string]bool{}
	for _, f := range s.order {
		if _, ok := out[f]; ok {
			ordered = append(ordered, f)
			seen[f] = true
		}
	}
	for k := range out {
		if !seen[k] {
			ordered = append(ordered, k)
		}
	}
	var sb strings.Builder
	sb.WriteByte('{')
	for i, k := range ordered {
		if i > 0 {
			sb.WriteByte(',')
		}
		kb, _ := json.Marshal(k)
		sb.Write(kb)
		sb.WriteByte(':')
		sb.Write(out[k])
	}
	sb.WriteByte('}')
	return sb.String(), true
}

func (s *Sorter) sortKV(line string) (string, bool) {
	pairs := strings.Fields(line)
	kvMap := map[string]string{}
	kvOrder := []string{}
	for _, p := range pairs {
		if idx := strings.IndexByte(p, '='); idx > 0 {
			k := p[:idx]
			v := p[idx+1:]
			if _, exists := kvMap[k]; !exists {
				kvOrder = append(kvOrder, k)
			}
			kvMap[k] = v
		} else {
			return "", false
		}
	}
	if len(kvMap) == 0 {
		return "", false
	}
	seen := map[string]bool{}
	result := []string{}
	for _, f := range s.order {
		if v, ok := kvMap[f]; ok {
			result = append(result, f+"="+v)
			seen[f] = true
		}
	}
	for _, k := range kvOrder {
		if !seen[k] {
			result = append(result, k+"="+kvMap[k])
		}
	}
	return strings.Join(result, " "), true
}
