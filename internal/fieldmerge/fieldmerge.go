// Package fieldmerge combines two log fields into a single new field.
package fieldmerge

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Merger combines two source fields into a destination field using a separator.
type Merger struct {
	srcA      string
	srcB      string
	dst       string
	sep       string
	merged    int
	skipped   int
}

// New returns a Merger that concatenates srcA and srcB into dst using sep.
// Returns an error if any field name is empty or dst duplicates a source.
func New(srcA, srcB, dst, sep string) (*Merger, error) {
	if srcA == "" {
		return nil, fmt.Errorf("fieldmerge: srcA must not be empty")
	}
	if srcB == "" {
		return nil, fmt.Errorf("fieldmerge: srcB must not be empty")
	}
	if dst == "" {
		return nil, fmt.Errorf("fieldmerge: dst must not be empty")
	}
	if dst == srcA || dst == srcB {
		return nil, fmt.Errorf("fieldmerge: dst %q must differ from source fields", dst)
	}
	return &Merger{srcA: srcA, srcB: srcB, dst: dst, sep: sep}, nil
}

// Apply merges the two source fields and returns the rewritten line.
// If either source field is absent the line is returned unchanged.
func (m *Merger) Apply(line string) string {
	line = strings.TrimRight(line, "\n")
	if line == "" {
		return line
	}
	var out string
	if strings.HasPrefix(line, "{") {
		out = m.applyJSON(line)
	} else {
		out = m.applyKV(line)
	}
	return out
}

func (m *Merger) applyJSON(line string) string {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		m.skipped++
		return line
	}
	va, okA := obj[m.srcA]
	vb, okB := obj[m.srcB]
	if !okA || !okB {
		m.skipped++
		return line
	}
	obj[m.dst] = fmt.Sprintf("%v%s%v", va, m.sep, vb)
	b, err := json.Marshal(obj)
	if err != nil {
		m.skipped++
		return line
	}
	m.merged++
	return string(b)
}

func (m *Merger) applyKV(line string) string {
	pairs := strings.Fields(line)
	vals := make(map[string]string, len(pairs))
	for _, p := range pairs {
		if k, v, ok := strings.Cut(p, "="); ok {
			vals[k] = v
		}
	}
	va, okA := vals[m.srcA]
	vb, okB := vals[m.srcB]
	if !okA || !okB {
		m.skipped++
		return line
	}
	return line + " " + m.dst + "=" + va + m.sep + vb
}

// Merged returns the number of lines where the merge was applied.
func (m *Merger) Merged() int { return m.merged }

// Skipped returns the number of lines where a source field was absent.
func (m *Merger) Skipped() int { return m.skipped }
