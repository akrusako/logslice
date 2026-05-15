// Package fieldcast coerces a named log field to a target type (string, int,
// float, bool) and rewrites the line with the converted value.
package fieldcast

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/yourorg/logslice/internal/fieldextract"
)

// Type enumerates the supported cast targets.
type Type string

const (
	TypeString Type = "string"
	TypeInt    Type = "int"
	TypeFloat  Type = "float"
	TypeBool   Type = "bool"
)

// ParseType converts a string to a Type or returns an error.
func ParseType(s string) (Type, error) {
	switch Type(strings.ToLower(s)) {
	case TypeString, TypeInt, TypeFloat, TypeBool:
		return Type(strings.ToLower(s)), nil
	}
	return "", fmt.Errorf("fieldcast: unknown type %q; want string|int|float|bool", s)
}

// Caster rewrites a single field value to the requested type.
type Caster struct {
	field   string
	to      Type
	casted  int64
	failed  int64
}

// New creates a Caster that converts field to the given type.
func New(field string, to Type) (*Caster, error) {
	if field == "" {
		return nil, fmt.Errorf("fieldcast: field name must not be empty")
	}
	return &Caster{field: field, to: to}, nil
}

// Apply rewrites the field value in line and returns the modified line.
// Lines where the field is absent are returned unchanged.
func (c *Caster) Apply(line string) string {
	fields := fieldextract.Extract(line)
	raw, ok := fields[c.field]
	if !ok {
		return line
	}
	converted, err := convert(raw, c.to)
	if err != nil {
		c.failed++
		return line
	}
	c.casted++
	return rewrite(line, c.field, converted)
}

// Casted returns the number of successful conversions.
func (c *Caster) Casted() int64 { return c.casted }

// Failed returns the number of lines where conversion failed.
func (c *Caster) Failed() int64 { return c.failed }

func convert(raw string, to Type) (string, error) {
	switch to {
	case TypeString:
		return raw, nil
	case TypeInt:
		f, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			return "", err
		}
		return strconv.FormatInt(int64(f), 10), nil
	case TypeFloat:
		f, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			return "", err
		}
		return strconv.FormatFloat(f, 'f', -1, 64), nil
	case TypeBool:
		b, err := strconv.ParseBool(raw)
		if err != nil {
			// treat "1"/"0" explicitly
			if raw == "1" {
				return "true", nil
			} else if raw == "0" {
				return "false", nil
			}
			return "", err
		}
		return strconv.FormatBool(b), nil
	}
	return "", fmt.Errorf("unknown type")
}

// rewrite replaces key=oldVal or "key":"oldVal" with key=newVal / "key":newVal.
func rewrite(line, field, newVal string) string {
	// Try JSON style first.
	if strings.HasPrefix(strings.TrimSpace(line), "{") {
		old := fmt.Sprintf(`"%s":`, field)
		idx := strings.Index(line, old)
		if idx >= 0 {
			after := line[idx+len(old):]
			// find end of value token
			end := strings.IndexAny(after, ",}")
			if end < 0 {
				end = len(after)
			}
			return line[:idx+len(old)] + newVal + after[end:]
		}
	}
	// KV style.
	kvKey := field + "="
	idx := strings.Index(line, kvKey)
	if idx < 0 {
		return line
	}
	after := line[idx+len(kvKey):]
	end := strings.IndexByte(after, ' ')
	if end < 0 {
		end = len(after)
	}
	return line[:idx+len(kvKey)] + newVal + after[end:]
}
