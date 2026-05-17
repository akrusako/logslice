// Package fieldpad left- or right-pads a named field value to a fixed width.
package fieldpad

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Direction controls which side padding is applied to.
type Direction int

const (
	Left  Direction = iota // pad on the left (right-align)
	Right                  // pad on the right (left-align)
)

// Padder pads a single log field to a target width.
type Padder struct {
	field  string
	width  int
	char   rune
	dir    Direction
	padded int
}

// New returns a Padder for field, targeting width columns using padChar.
// dir selects Left or Right padding. Returns an error for invalid arguments.
func New(field string, width int, padChar rune, dir Direction) (*Padder, error) {
	if field == "" {
		return nil, fmt.Errorf("fieldpad: field name must not be empty")
	}
	if width <= 0 {
		return nil, fmt.Errorf("fieldpad: width must be positive, got %d", width)
	}
	if padChar == 0 {
		padChar = ' '
	}
	return &Padder{field: field, width: width, char: padChar, dir: dir}, nil
}

// Apply pads the target field in line and returns the result.
// Lines that do not contain the field are returned unchanged.
func (p *Padder) Apply(line string) string {
	if strings.HasPrefix(strings.TrimSpace(line), "{") {
		return p.applyJSON(line)
	}
	return p.applyKV(line)
}

// PaddedCount returns the number of lines on which padding was applied.
func (p *Padder) PaddedCount() int { return p.padded }

func (p *Padder) pad(s string) string {
	if len(s) >= p.width {
		return s
	}
	fill := strings.Repeat(string(p.char), p.width-len(s))
	if p.dir == Left {
		return fill + s
	}
	return s + fill
}

func (p *Padder) applyJSON(line string) string {
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		return line
	}
	v, ok := m[p.field]
	if !ok {
		return line
	}
	original := fmt.Sprintf("%v", v)
	padded := p.pad(original)
	if padded == original {
		return line
	}
	m[p.field] = padded
	b, err := json.Marshal(m)
	if err != nil {
		return line
	}
	p.padded++
	return string(b)
}

func (p *Padder) applyKV(line string) string {
	parts := strings.Fields(line)
	changed := false
	for i, part := range parts {
		if !strings.Contains(part, "=") {
			continue
		}
		idx := strings.IndexByte(part, '=')
		key := part[:idx]
		val := part[idx+1:]
		if key != p.field {
			continue
		}
		padded := p.pad(val)
		if padded != val {
			parts[i] = key + "=" + padded
			changed = true
		}
	}
	if !changed {
		return line
	}
	p.padded++
	return strings.Join(parts, " ")
}
