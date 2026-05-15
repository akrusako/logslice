// Package fieldclip clamps numeric field values to a [min, max] range,
// rewriting the field in-place for both KV and JSON log lines.
package fieldclip

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// Clipper clamps a named numeric field to [Min, Max].
type Clipper struct {
	field   string
	min     float64
	max     float64
	clipped int
}

// New returns a Clipper for the given field and bounds.
// Returns an error if field is empty or min > max.
func New(field string, min, max float64) (*Clipper, error) {
	if field == "" {
		return nil, fmt.Errorf("fieldclip: field name must not be empty")
	}
	if min > max {
		return nil, fmt.Errorf("fieldclip: min %.6g > max %.6g", min, max)
	}
	return &Clipper{field: field, min: min, max: max}, nil
}

// Apply clamps the field value in line and returns the rewritten line.
func (c *Clipper) Apply(line string) string {
	if out, ok := c.applyJSON(line); ok {
		return out
	}
	if out, ok := c.applyKV(line); ok {
		return out
	}
	return line
}

// Clipped returns the total number of values that were clamped.
func (c *Clipper) Clipped() int { return c.clipped }

func (c *Clipper) clamp(v float64) (float64, bool) {
	if v < c.min {
		return c.min, true
	}
	if v > c.max {
		return c.max, true
	}
	return v, false
}

func (c *Clipper) applyJSON(line string) (string, bool) {
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		return "", false
	}
	v, ok := m[c.field]
	if !ok {
		return line, true
	}
	num, ok2 := toFloat(v)
	if !ok2 {
		return line, true
	}
	clamped, changed := c.clamp(num)
	if changed {
		c.clipped++
		m[c.field] = clamped
	}
	b, err := json.Marshal(m)
	if err != nil {
		return line, true
	}
	return string(b), true
}

func (c *Clipper) applyKV(line string) (string, bool) {
	prefix := c.field + "="
	idx := strings.Index(line, prefix)
	if idx == -1 {
		return "", false
	}
	rest := line[idx+len(prefix):]
	end := strings.IndexByte(rest, ' ')
	var raw string
	if end == -1 {
		raw = rest
	} else {
		raw = rest[:end]
	}
	num, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return "", false
	}
	clamped, changed := c.clamp(num)
	if !changed {
		return line, true
	}
	c.clipped++
	newVal := strconv.FormatFloat(clamped, 'f', -1, 64)
	return line[:idx+len(prefix)] + newVal + line[idx+len(prefix)+len(raw):], true
}

func toFloat(v interface{}) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case int:
		return float64(n), true
	}
	return 0, false
}
