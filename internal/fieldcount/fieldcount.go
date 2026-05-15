// Package fieldcount counts the number of fields present in each log line
// and allows filtering lines by a minimum or maximum field count.
package fieldcount

import (
	"fmt"

	"github.com/yourorg/logslice/internal/fieldextract"
)

// Filter holds configuration for field-count filtering.
type Filter struct {
	min     int
	max     int // -1 means no upper bound
	allowed int
	dropped int
}

// New creates a Filter that keeps lines whose field count is in [min, max].
// Set max to -1 to disable the upper bound.
func New(min, max int) (*Filter, error) {
	if min < 0 {
		return nil, fmt.Errorf("fieldcount: min must be >= 0, got %d", min)
	}
	if max != -1 && max < min {
		return nil, fmt.Errorf("fieldcount: max (%d) must be >= min (%d) or -1", max, min)
	}
	return &Filter{min: min, max: max}, nil
}

// Allow returns true when the number of fields in line falls within [min, max].
func (f *Filter) Allow(line string) bool {
	fields := fieldextract.Extract(line)
	n := len(fields)
	ok := n >= f.min && (f.max == -1 || n <= f.max)
	if ok {
		f.allowed++
	} else {
		f.dropped++
	}
	return ok
}

// Allowed returns the total number of lines that passed the filter.
func (f *Filter) Allowed() int { return f.allowed }

// Dropped returns the total number of lines that were rejected.
func (f *Filter) Dropped() int { return f.dropped }

// Count returns the number of fields detected in line without updating counters.
func Count(line string) int {
	return len(fieldextract.Extract(line))
}
