// Package multiline merges continuation lines into a single logical log entry.
// Lines that match the start pattern begin a new entry; lines that do not match
// are appended to the current entry until the next start pattern is seen.
package multiline

import (
	"fmt"
	"regexp"
	"strings"
)

// Merger accumulates continuation lines and emits complete entries.
type Merger struct {
	startRe  *regexp.Regexp
	joinSep  string
	pending  strings.Builder
	merged   uint64
	flushed  uint64
}

// New creates a Merger. startPattern is a regular expression that identifies
// the first line of a new log entry. joinSep is the string used to join
// continuation lines (commonly " " or "\n").
func New(startPattern, joinSep string) (*Merger, error) {
	if startPattern == "" {
		return nil, fmt.Errorf("multiline: startPattern must not be empty")
	}
	re, err := regexp.Compile(startPattern)
	if err != nil {
		return nil, fmt.Errorf("multiline: invalid pattern: %w", err)
	}
	return &Merger{startRe: re, joinSep: joinSep}, nil
}

// Feed accepts a raw line. It returns a complete entry string and true when
// the previous entry has been flushed, or "", false when the line was buffered.
func (m *Merger) Feed(line string) (string, bool) {
	if m.startRe.MatchString(line) {
		if m.pending.Len() == 0 {
			m.pending.WriteString(line)
			return "", false
		}
		entry := m.pending.String()
		m.pending.Reset()
		m.pending.WriteString(line)
		m.flushed++
		return entry, true
	}
	// continuation line
	if m.pending.Len() > 0 {
		m.pending.WriteString(m.joinSep)
		m.merged++
	}
	m.pending.WriteString(line)
	return "", false
}

// Flush returns any buffered entry that has not yet been emitted.
// Call after the input stream is exhausted.
func (m *Merger) Flush() (string, bool) {
	if m.pending.Len() == 0 {
		return "", false
	}
	entry := m.pending.String()
	m.pending.Reset()
	m.flushed++
	return entry, true
}

// Merged returns the number of continuation lines that were merged.
func (m *Merger) Merged() uint64 { return m.merged }

// Flushed returns the number of complete entries emitted.
func (m *Merger) Flushed() uint64 { return m.flushed }
