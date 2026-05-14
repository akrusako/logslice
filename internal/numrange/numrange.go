// Package numrange filters log lines by checking whether a numeric field
// falls within a specified inclusive [min, max] range.
package numrange

import (
	"fmt"
	"strconv"

	"github.com/user/logslice/internal/fieldextract"
)

// Filter drops lines where the named numeric field is outside [Min, Max].
type Filter struct {
	field   string
	min     float64
	max     float64
	passed  int64
	dropped int64
}

// New creates a Filter for the given field and inclusive numeric bounds.
// Returns an error if field is empty or min > max.
func New(field string, min, max float64) (*Filter, error) {
	if field == "" {
		return nil, fmt.Errorf("numrange: field name must not be empty")
	}
	if min > max {
		return nil, fmt.Errorf("numrange: min %.4g > max %.4g", min, max)
	}
	return &Filter{field: field, min: min, max: max}, nil
}

// Allow returns true when the field value is within [min, max].
// Lines that do not contain the field are passed through unchanged.
func (f *Filter) Allow(line string) bool {
	fields := fieldextract.Extract(line)
	raw, ok := fields[f.field]
	if !ok {
		f.passed++
		return true
	}
	v, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		f.dropped++
		return false
	}
	if v < f.min || v > f.max {
		f.dropped++
		return false
	}
	f.passed++
	return true
}

// Passed returns the number of lines allowed through.
func (f *Filter) Passed() int64 { return f.passed }

// Dropped returns the number of lines rejected.
func (f *Filter) Dropped() int64 { return f.dropped }
