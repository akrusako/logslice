// Package tail provides utilities for reading the last N lines
// of a log file efficiently using a pre-built line index.
package tail

import (
	"fmt"
	"io"
	"os"

	"github.com/yourorg/logslice/internal/lineindex"
)

// Reader holds state for tail operations on a file.
type Reader struct {
	path  string
	index lineindex.Index
}

// New opens the file at path, builds a line index, and returns a Reader.
func New(path string) (*Reader, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("tail: open %q: %w", path, err)
	}
	defer f.Close()

	idx, err := lineindex.Build(f)
	if err != nil {
		return nil, fmt.Errorf("tail: build index: %w", err)
	}
	return &Reader{path: path, index: idx}, nil
}

// Lines returns the last n lines of the file.
// If n is greater than the total number of lines, all lines are returned.
func (r *Reader) Lines(n int) ([]string, error) {
	total := r.index.Len()
	if n <= 0 {
		return nil, nil
	}
	if n > total {
		n = total
	}

	startLine := total - n // 0-based

	offset, err := r.index.OffsetForLine(startLine)
	if err != nil {
		return nil, fmt.Errorf("tail: offset lookup: %w", err)
	}

	f, err := os.Open(r.path)
	if err != nil {
		return nil, fmt.Errorf("tail: open %q: %w", r.path, err)
	}
	defer f.Close()

	if _, err := f.Seek(offset, io.SeekStart); err != nil {
		return nil, fmt.Errorf("tail: seek: %w", err)
	}

	return readLines(f, n)
}

// readLines reads up to max lines from r.
func readLines(r io.Reader, max int) ([]string, error) {
	buf := make([]byte, 0, 4096)
	all, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	buf = all

	var lines []string
	start := 0
	for i, b := range buf {
		if b == '\n' {
			lines = append(lines, string(buf[start:i]))
			start = i + 1
			if len(lines) == max {
				break
			}
		}
	}
	if start < len(buf) && len(lines) < max {
		lines = append(lines, string(buf[start:]))
	}
	return lines, nil
}
