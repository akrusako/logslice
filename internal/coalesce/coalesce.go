// Package coalesce provides a processor that merges multi-line log entries
// into a single logical line based on a continuation pattern.
package coalesce

import (
	"regexp"
	"strings"
)

// Coalescer merges continuation lines into the preceding anchor line.
type Coalescer struct {
	continuation *regexp.Regexp
	pending      strings.Builder
	hasPending   bool
	joiner       string
	merged       int
}

// New creates a Coalescer that treats lines matching continuationPattern as
// continuations of the previous line. joiner is inserted between segments
// (use " " or "\n" depending on desired output).
func New(continuationPattern, joiner string) (*Coalescer, error) {
	re, err := regexp.Compile(continuationPattern)
	if err != nil {
		return nil, err
	}
	return &Coalescer{
		continuation: re,
		joiner:       joiner,
	}, nil
}

// Process feeds a line into the coalescer. It returns the flushed line and
// true when a complete logical entry is ready, or "", false when the line was
// buffered as a continuation.
func (c *Coalescer) Process(line string) (string, bool) {
	if c.continuation.MatchString(line) {
		// Continuation: append to pending buffer.
		if c.hasPending {
			c.pending.WriteString(c.joiner)
		}
		c.pending.WriteString(line)
		c.hasPending = true
		c.merged++
		return "", false
	}

	// New anchor line: flush any pending entry first.
	if c.hasPending {
		flushed := c.pending.String()
		c.pending.Reset()
		c.hasPending = false
		// Start buffering the new anchor.
		c.pending.WriteString(line)
		c.hasPending = true
		return flushed, true
	}

	// No pending buffer: start a new anchor.
	c.pending.WriteString(line)
	c.hasPending = true
	return "", false
}

// Flush returns any buffered line that has not yet been emitted. Call this
// after the last line of input to drain the internal buffer.
func (c *Coalescer) Flush() (string, bool) {
	if !c.hasPending {
		return "", false
	}
	result := c.pending.String()
	c.pending.Reset()
	c.hasPending = false
	return result, true
}

// MergedCount returns the total number of continuation lines absorbed.
func (c *Coalescer) MergedCount() int { return c.merged }
