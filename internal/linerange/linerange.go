package linerange

import (
	"fmt"
	"strconv"
	"strings"
)

// Range represents an inclusive line number range [Start, End].
// Line numbers are 1-based. A zero value for End means "no upper bound".
type Range struct {
	Start int
	End   int // 0 means unbounded
}

// Contains reports whether the given 1-based line number falls within the range.
func (r Range) Contains(line int) bool {
	if line < r.Start {
		return false
	}
	if r.End != 0 && line > r.End {
		return false
	}
	return true
}

// String returns a human-readable representation of the range.
func (r Range) String() string {
	if r.End == 0 {
		return fmt.Sprintf("%d:", r.Start)
	}
	return fmt.Sprintf("%d:%d", r.Start, r.End)
}

// Parse parses a line range string of the form "N", "N:", "N:M", or ":M".
// Line numbers are 1-based. Returns an error if the format is invalid or
// if Start > End (when both are specified).
func Parse(s string) (Range, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return Range{}, fmt.Errorf("linerange: empty range string")
	}

	if !strings.Contains(s, ":") {
		// Single line number
		n, err := parsePositive(s)
		if err != nil {
			return Range{}, fmt.Errorf("linerange: %w", err)
		}
		return Range{Start: n, End: n}, nil
	}

	parts := strings.SplitN(s, ":", 2)
	var start, end int
	var err error

	if parts[0] == "" {
		start = 1
	} else {
		start, err = parsePositive(parts[0])
		if err != nil {
			return Range{}, fmt.Errorf("linerange: start: %w", err)
		}
	}

	if parts[1] == "" {
		end = 0 // unbounded
	} else {
		end, err = parsePositive(parts[1])
		if err != nil {
			return Range{}, fmt.Errorf("linerange: end: %w", err)
		}
		if end < start {
			return Range{}, fmt.Errorf("linerange: end %d < start %d", end, start)
		}
	}

	return Range{Start: start, End: end}, nil
}

func parsePositive(s string) (int, error) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("invalid number %q", s)
	}
	if n < 1 {
		return 0, fmt.Errorf("line number must be >= 1, got %d", n)
	}
	return n, nil
}
