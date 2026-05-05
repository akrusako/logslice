// Package logscanner provides line-by-line log scanning with timestamp
// extraction and time-range filtering capabilities.
package logscanner

import (
	"bufio"
	"io"
	"time"

	"github.com/yourorg/logslice/internal/timeparse"
)

// Line represents a single parsed log line.
type Line struct {
	Raw       string
	Timestamp time.Time
	HasTime   bool
	Number    int64
}

// Scanner reads log lines and optionally filters by a time range.
type Scanner struct {
	parser  *timeparse.Parser
	reader  *bufio.Reader
	after   time.Time
	before  time.Time
	filter  bool
	current Line
	lineNum int64
	err     error
}

// Options configures a Scanner.
type Options struct {
	// After filters out lines with timestamps before this time (zero = no lower bound).
	After time.Time
	// Before filters out lines with timestamps after this time (zero = no upper bound).
	Before time.Time
}

// New creates a Scanner that reads from r.
func New(r io.Reader, opts Options) *Scanner {
	return &Scanner{
		parser: timeparse.NewParser(),
		reader: bufio.NewReaderSize(r, 64*1024),
		after:  opts.After,
		before: opts.Before,
		filter: !opts.After.IsZero() || !opts.Before.IsZero(),
	}
}

// Scan advances to the next matching line. Returns false when done or on error.
func (s *Scanner) Scan() bool {
	for {
		raw, err := s.reader.ReadString('\n')
		if len(raw) == 0 {
			if err == io.EOF {
				return false
			}
			if err != nil {
				s.err = err
				return false
			}
		}
		s.lineNum++
		line := Line{
			Raw:    raw,
			Number: s.lineNum,
		}
		if ts, ok := s.parser.ParseAny(raw); ok {
			line.Timestamp = ts
			line.HasTime = true
		}
		if s.filter && line.HasTime {
			if !s.after.IsZero() && line.Timestamp.Before(s.after) {
				continue
			}
			if !s.before.IsZero() && line.Timestamp.After(s.before) {
				continue
			}
		}
		s.current = line
		return true
	}
}

// Line returns the current log line.
func (s *Scanner) Line() Line { return s.current }

// Err returns the first non-EOF error encountered.
func (s *Scanner) Err() error { return s.err }
