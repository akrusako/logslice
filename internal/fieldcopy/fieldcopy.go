// Package fieldcopy duplicates a field value under a new key in structured log lines.
package fieldcopy

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copier duplicates a named field to a new key in each log line.
type Copier struct {
	src    string
	dst    string
	copied int
}

// New returns a Copier that copies the value of src into dst.
// Both src and dst must be non-empty and distinct.
func New(src, dst string) (*Copier, error) {
	if src == "" {
		return nil, fmt.Errorf("fieldcopy: src field must not be empty")
	}
	if dst == "" {
		return nil, fmt.Errorf("fieldcopy: dst field must not be empty")
	}
	if src == dst {
		return nil, fmt.Errorf("fieldcopy: src and dst must differ")
	}
	return &Copier{src: src, dst: dst}, nil
}

// Apply returns the line with the src field value duplicated under dst.
// Lines where src is absent are returned unchanged.
func (c *Copier) Apply(line string) string {
	if out, ok := applyJSON(line, c.src, c.dst); ok {
		c.copied++
		return out
	}
	if out, ok := applyKV(line, c.src, c.dst); ok {
		c.copied++
		return out
	}
	return line
}

// Copied returns the number of lines where a copy was performed.
func (c *Copier) Copied() int { return c.copied }

func applyJSON(line, src, dst string) (string, bool) {
	if !strings.HasPrefix(strings.TrimSpace(line), "{") {
		return "", false
	}
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		return "", false
	}
	val, ok := m[src]
	if !ok {
		return "", false
	}
	m[dst] = val
	b, err := json.Marshal(m)
	if err != nil {
		return "", false
	}
	return string(b), true
}

func applyKV(line, src, dst string) (string, bool) {
	prefix := src + "="
	idx := strings.Index(line, prefix)
	if idx == -1 {
		return "", false
	}
	rest := line[idx+len(prefix):]
	var val string
	if len(rest) > 0 && rest[0] == '"' {
		end := strings.Index(rest[1:], "\"")
		if end == -1 {
			return "", false
		}
		val = rest[:end+2]
	} else {
		end := strings.IndexByte(rest, ' ')
		if end == -1 {
			val = rest
		} else {
			val = rest[:end]
		}
	}
	return line + " " + dst + "=" + val, true
}
