// Package bsearch provides binary search helpers for locating log lines
// within a time range using a pre-built line index.
package bsearch

import (
	"io"
	"time"

	"github.com/yourorg/logslice/internal/lineindex"
	"github.com/yourorg/logslice/internal/timeparse"
)

// Finder locates the first and last line offsets that fall within a
// [start, end] time range using binary search over the line index.
type Finder struct {
	index  *lineindex.Index
	parser *timeparse.Parser
	r      io.ReaderAt
}

// New creates a Finder backed by the given index, reader, and time parser.
func New(idx *lineindex.Index, r io.ReaderAt, p *timeparse.Parser) *Finder {
	return &Finder{index: idx, r: r, parser: p}
}

// timestampAt reads the line at position i and attempts to parse its leading
// timestamp. It returns the zero Time and false when parsing fails.
func (f *Finder) timestampAt(i int) (time.Time, bool) {
	off, length, ok := f.index.OffsetForLine(i)
	if !ok {
		return time.Time{}, false
	}
	buf := make([]byte, length)
	n, err := f.r.ReadAt(buf, off)
	if err != nil && n == 0 {
		return time.Time{}, false
	}
	t, err := f.parser.ParseAny(string(buf[:n]))
	if err != nil {
		return time.Time{}, false
	}
	return t, true
}

// FindRange returns the inclusive [firstLine, lastLine] line numbers whose
// timestamps fall within [start, end]. Both values are -1 when no lines match.
func (f *Finder) FindRange(start, end time.Time) (first, last int) {
	n := f.index.Len()
	if n == 0 {
		return -1, -1
	}

	// Binary search for the first line with timestamp >= start.
	lo, hi := 0, n-1
	first = -1
	for lo <= hi {
		mid := (lo + hi) / 2
		t, ok := f.timestampAt(mid)
		if !ok || t.Before(start) {
			lo = mid + 1
		} else {
			first = mid
			hi = mid - 1
		}
	}
	if first == -1 {
		return -1, -1
	}

	// Binary search for the last line with timestamp <= end.
	lo, hi = first, n-1
	last = -1
	for lo <= hi {
		mid := (lo + hi) / 2
		t, ok := f.timestampAt(mid)
		if !ok || t.After(end) {
			hi = mid - 1
		} else {
			last = mid
			lo = mid + 1
		}
	}
	if last < first {
		return -1, -1
	}
	return first, last
}
