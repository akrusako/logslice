// Package fieldrename provides a processor that renames fields in structured
// log lines (JSON or key=value format) according to a mapping of old→new names.
package fieldrename

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Renamer renames fields in log lines.
type Renamer struct {
	mapping  map[string]string
	renamed  int64
	processed int64
}

// New creates a Renamer from a slice of "old=new" pair strings.
// Returns an error if any pair is malformed or either side is empty.
func New(pairs []string) (*Renamer, error) {
	mapping := make(map[string]string, len(pairs))
	for _, p := range pairs {
		parts := strings.SplitN(p, "=", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("fieldrename: invalid pair %q: must be old=new", p)
		}
		mapping[parts[0]] = parts[1]
	}
	return &Renamer{mapping: mapping}, nil
}

// Apply renames fields in line according to the configured mapping.
// Lines that cannot be parsed as JSON or KV are returned unchanged.
func (r *Renamer) Apply(line string) string {
	r.processed++
	if len(r.mapping) == 0 {
		return line
	}
	if result, ok := r.applyJSON(line); ok {
		r.renamed++
		return result
	}
	if result, ok := r.applyKV(line); ok {
		r.renamed++
		return result
	}
	return line
}

func (r *Renamer) applyJSON(line string) (string, bool) {
	if !strings.HasPrefix(strings.TrimSpace(line), "{") {
		return "", false
	}
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		return "", false
	}
	changed := false
	for old, newKey := range r.mapping {
		if v, exists := m[old]; exists {
			delete(m, old)
			m[newKey] = v
			changed = true
		}
	}
	if !changed {
		return "", false
	}
	b, err := json.Marshal(m)
	if err != nil {
		return "", false
	}
	return string(b), true
}

func (r *Renamer) applyKV(line string) (string, bool) {
	parts := strings.Fields(line)
	changed := false
	for i, part := range parts {
		idx := strings.IndexByte(part, '=')
		if idx < 1 {
			continue
		}
		key := part[:idx]
		if newKey, ok := r.mapping[key]; ok {
			parts[i] = newKey + part[idx:]
			changed = true
		}
	}
	if !changed {
		return "", false
	}
	return strings.Join(parts, " "), true
}

// Renamed returns the number of lines that had at least one field renamed.
func (r *Renamer) Renamed() int64 { return r.renamed }

// Processed returns the total number of lines processed.
func (r *Renamer) Processed() int64 { return r.processed }
