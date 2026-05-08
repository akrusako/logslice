// Package head provides a processor that passes through only the first N lines
// of a log stream, dropping all subsequent lines once the limit is reached.
package head

import "fmt"

// Processor limits output to the first N lines.
type Processor struct {
	max     int
	count   int
	dropped int
}

// New creates a Processor that passes through at most max lines.
// If max is zero or negative, all lines are passed through.
func New(max int) *Processor {
	return &Processor{max: max}
}

// Process returns the line unchanged if the head limit has not been reached,
// or an empty string once the limit is exceeded. A non-empty string signals
// the caller to emit the line; an empty string signals a drop.
func (p *Processor) Process(line string) (string, bool) {
	if p.max <= 0 {
		return line, true
	}
	if p.count < p.max {
		p.count++
		return line, true
	}
	p.dropped++
	return "", false
}

// Done returns true when the head limit has been reached and no further lines
// will be emitted. Callers may use this to short-circuit reading.
func (p *Processor) Done() bool {
	return p.max > 0 && p.count >= p.max
}

// Stats returns a human-readable summary of lines kept and dropped.
func (p *Processor) Stats() string {
	return fmt.Sprintf("head: kept %d, dropped %d", p.count, p.dropped)
}
