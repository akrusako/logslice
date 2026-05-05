package filter

import (
	"fmt"
	"strings"
)

// ParseFields parses a slice of "key=value" strings into a map.
// It returns an error if any entry is not in key=value form.
func ParseFields(pairs []string) (map[string]string, error) {
	if len(pairs) == 0 {
		return nil, nil
	}
	out := make(map[string]string, len(pairs))
	for _, p := range pairs {
		parts := strings.SplitN(p, "=", 2)
		if len(parts) != 2 || parts[0] == "" {
			return nil, fmt.Errorf("filter: invalid field expression %q (want key=value)", p)
		}
		out[parts[0]] = parts[1]
	}
	return out, nil
}

// NewFromArgs is a convenience constructor that parses raw "key=value" strings
// and an optional substring into a ready-to-use Filter.
func NewFromArgs(fieldPairs []string, substring string) (*Filter, error) {
	fields, err := ParseFields(fieldPairs)
	if err != nil {
		return nil, err
	}
	return New(Options{
		Fields:    fields,
		Substring: substring,
	}), nil
}
