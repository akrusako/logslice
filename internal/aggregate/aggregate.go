// Package aggregate provides field-based log line counting and grouping.
package aggregate

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/logslice/logslice/internal/fieldextract"
)

// Counter accumulates counts keyed by a field value.
type Counter struct {
	field  string
	counts map[string]int
	total  int
}

// New creates a Counter that groups lines by the given field name.
// If field is empty, every line is counted under the key "*".
func New(field string) *Counter {
	return &Counter{
		field:  field,
		counts: make(map[string]int),
	}
}

// Record ingests a single log line and increments the appropriate bucket.
func (c *Counter) Record(line string) {
	c.total++
	key := "*"
	if c.field != "" {
		fields := fieldextract.Extract(line)
		if v, ok := fields[c.field]; ok {
			key = v
		} else {
			key = "(missing)"
		}
	}
	c.counts[key]++
}

// Total returns the number of lines recorded.
func (c *Counter) Total() int { return c.total }

// Counts returns a copy of the internal bucket map.
func (c *Counter) Counts() map[string]int {
	out := make(map[string]int, len(c.counts))
	for k, v := range c.counts {
		out[k] = v
	}
	return out
}

// WriteSummary writes a sorted summary table to w.
func (c *Counter) WriteSummary(w io.Writer) error {
	keys := make([]string, 0, len(c.counts))
	for k := range c.counts {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		if c.counts[keys[i]] != c.counts[keys[j]] {
			return c.counts[keys[i]] > c.counts[keys[j]]
		}
		return keys[i] < keys[j]
	})
	var sb strings.Builder
	for _, k := range keys {
		sb.WriteString(fmt.Sprintf("%6d\t%s\n", c.counts[k], k))
	}
	_, err := io.WriteString(w, sb.String())
	return err
}
