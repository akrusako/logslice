// Package linenum injects a line number prefix into each log line.
package linenum

import "fmt"

// Injector prepends a sequential line number to each processed line.
type Injector struct {
	start   int
	current int
	format  string
	count   int
}

// New creates an Injector that begins numbering at start.
// The format string must contain exactly one %d verb used for the line number.
// If format is empty, the default "[%d] " is used.
func New(start int, format string) *Injector {
	if format == "" {
		format = "[%d] "
	}
	return &Injector{
		start:   start,
		current: start,
		format:  format,
	}
}

// Process prepends the current line number to line and increments the counter.
// Empty lines are returned unchanged and do not advance the counter.
func (inj *Injector) Process(line string) string {
	if line == "" {
		return line
	}
	prefix := fmt.Sprintf(inj.format, inj.current)
	inj.current++
	inj.count++
	return prefix + line
}

// Reset resets the counter back to the start value provided at construction.
func (inj *Injector) Reset() {
	inj.current = inj.start
	inj.count = 0
}

// Count returns the number of lines that have been numbered so far.
func (inj *Injector) Count() int {
	return inj.count
}

// Current returns the next line number that will be assigned.
func (inj *Injector) Current() int {
	return inj.current
}
