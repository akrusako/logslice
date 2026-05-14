// Package numfield extracts a numeric value from a named field in structured
// log lines and filters lines based on a configurable threshold comparison.
package numfield

import (
	"fmt"
	"strconv"

	"github.com/user/logslice/internal/fieldextract"
)

// Op represents a numeric comparison operator.
type Op int

const (
	OpGT  Op = iota // greater than
	OpGTE            // greater than or equal
	OpLT             // less than
	OpLTE            // less than or equal
	OpEQ             // equal
)

// ParseOp converts a string operator symbol to an Op.
func ParseOp(s string) (Op, error) {
	switch s {
	case ">":
		return OpGT, nil
	case ">=":
		return OpGTE, nil
	case "<":
		return OpLT, nil
	case "<=":
		return OpLTE, nil
	case "=", "==":
		return OpEQ, nil
	}
	return 0, fmt.Errorf("numfield: unknown operator %q", s)
}

// Filter passes log lines whose named numeric field satisfies the threshold.
type Filter struct {
	field     string
	op        Op
	threshold float64
	matched   int64
	dropped   int64
}

// New creates a Filter for the given field, operator, and threshold.
func New(field string, op Op, threshold float64) (*Filter, error) {
	if field == "" {
		return nil, fmt.Errorf("numfield: field name must not be empty")
	}
	return &Filter{field: field, op: op, threshold: threshold}, nil
}

// Allow returns true if the line's field value satisfies the comparison.
// Lines that do not contain the field are passed through unchanged.
func (f *Filter) Allow(line string) bool {
	fields := fieldextract.Extract(line)
	raw, ok := fields[f.field]
	if !ok {
		// no field present — pass through
		f.matched++
		return true
	}
	v, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		f.matched++
		return true
	}
	var pass bool
	switch f.op {
	case OpGT:
		pass = v > f.threshold
	case OpGTE:
		pass = v >= f.threshold
	case OpLT:
		pass = v < f.threshold
	case OpLTE:
		pass = v <= f.threshold
	case OpEQ:
		pass = v == f.threshold
	}
	if pass {
		f.matched++
	} else {
		f.dropped++
	}
	return pass
}

// Matched returns the number of lines that passed the filter.
func (f *Filter) Matched() int64 { return f.matched }

// Dropped returns the number of lines that were rejected by the filter.
func (f *Filter) Dropped() int64 { return f.dropped }
