// Package numformat rewrites numeric fields in log lines, applying formatting
// such as rounding, fixed decimal places, or SI suffixes (k, M, G).
package numformat

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Style controls how a numeric value is rendered.
type Style int

const (
	StyleFixed  Style = iota // fixed decimal places
	StyleSI                  // SI suffix (k, M, G, T)
	StyleRound               // round to nearest integer
)

// Formatter rewrites a named numeric field in each log line.
type Formatter struct {
	field   string
	style   Style
	prec    int
	formatted int64
	dropped   int64
}

// New creates a Formatter targeting field, using the given Style and precision
// (decimal places; ignored for StyleRound and StyleSI).
func New(field string, style Style, prec int) (*Formatter, error) {
	if field == "" {
		return nil, fmt.Errorf("numformat: field name must not be empty")
	}
	if prec < 0 {
		prec = 0
	}
	return &Formatter{field: field, style: style, prec: prec}, nil
}

// Apply rewrites the target field in line and returns the result.
// Lines where the field is absent or non-numeric are returned unchanged.
func (f *Formatter) Apply(line string) string {
	if strings.HasPrefix(strings.TrimSpace(line), "{") {
		out, ok := f.applyJSON(line)
		if ok {
			f.formatted++
			return out
		}
	} else {
		out, ok := f.applyKV(line)
		if ok {
			f.formatted++
			return out
		}
	}
	f.dropped++
	return line
}

func (f *Formatter) applyJSON(line string) (string, bool) {
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		return "", false
	}
	v, ok := m[f.field]
	if !ok {
		return line, true // field absent — pass through, not an error
	}
	num, ok := toFloat(v)
	if !ok {
		return "", false
	}
	m[f.field] = f.render(num)
	b, err := json.Marshal(m)
	if err != nil {
		return "", false
	}
	return string(b), true
}

func (f *Formatter) applyKV(line string) (string, bool) {
	parts := strings.Fields(line)
	found := false
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
		num, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return "", false
		}
		parts[i] = key + "=" + f.render(num)
		found = true
	}
	if !found {
		return line, true
	}
	return strings.Join(parts, " "), true
}

func (f *Formatter) render(v float64) string {
	switch f.style {
	case StyleSI:
		return siFormat(v)
	case StyleRound:
		return strconv.FormatInt(int64(math.Round(v)), 10)
	default:
		return strconv.FormatFloat(v, 'f', f.prec, 64)
	}
}

func siFormat(v float64) string {
	abs := math.Abs(v)
	switch {
	case abs >= 1e12:
		return strconv.FormatFloat(v/1e12, 'f', 2, 64) + "T"
	case abs >= 1e9:
		return strconv.FormatFloat(v/1e9, 'f', 2, 64) + "G"
	case abs >= 1e6:
		return strconv.FormatFloat(v/1e6, 'f', 2, 64) + "M"
	case abs >= 1e3:
		return strconv.FormatFloat(v/1e3, 'f', 2, 64) + "k"
	default:
		return strconv.FormatFloat(v, 'f', 2, 64)
	}
}

func toFloat(v interface{}) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case json.Number:
		f, err := n.Float64()
		return f, err == nil
	}
	return 0, false
}

// Formatted returns the number of lines where the field was successfully rewritten.
func (f *Formatter) Formatted() int64 { return f.formatted }

// Dropped returns the number of lines skipped due to parse errors.
func (f *Formatter) Dropped() int64 { return f.dropped }
