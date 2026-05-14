// Package labelinject attaches static key=value labels to every log line.
// Labels are appended as structured fields, preserving the original format
// (JSON or key=value). Plain-text lines receive a key=value suffix.
package labelinject

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Injector holds a set of static labels and injects them into log lines.
type Injector struct {
	labels  map[string]string
	order   []string // insertion order for deterministic output
	injected int
}

// New creates an Injector from a slice of "key=value" strings.
// Returns an error if any entry is malformed or the key is empty.
func New(pairs []string) (*Injector, error) {
	inj := &Injector{
		labels: make(map[string]string, len(pairs)),
		order:  make([]string, 0, len(pairs)),
	}
	for _, p := range pairs {
		k, v, ok := strings.Cut(p, "=")
		if !ok || k == "" {
			return nil, fmt.Errorf("labelinject: invalid label %q: must be key=value", p)
		}
		if _, dup := inj.labels[k]; !dup {
			inj.order = append(inj.order, k)
		}
		inj.labels[k] = v
	}
	return inj, nil
}

// Apply injects the configured labels into line and returns the result.
func (inj *Injector) Apply(line string) string {
	if len(inj.labels) == 0 {
		return line
	}
	result := inj.tryJSON(line)
	if result == "" {
		result = inj.appendKV(line)
	}
	inj.injected++
	return result
}

// Injected returns the total number of lines that had labels applied.
func (inj *Injector) Injected() int { return inj.injected }

func (inj *Injector) tryJSON(line string) string {
	trimmed := strings.TrimSpace(line)
	if len(trimmed) < 2 || trimmed[0] != '{' {
		return ""
	}
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(trimmed), &m); err != nil {
		return ""
	}
	for _, k := range inj.order {
		m[k] = inj.labels[k]
	}
	b, err := json.Marshal(m)
	if err != nil {
		return ""
	}
	return string(b)
}

func (inj *Injector) appendKV(line string) string {
	var sb strings.Builder
	sb.WriteString(line)
	for _, k := range inj.order {
		sb.WriteByte(' ')
		sb.WriteString(k)
		sb.WriteByte('=')
		sb.WriteString(inj.labels[k])
	}
	return sb.String()
}
