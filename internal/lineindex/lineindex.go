// Package lineindex provides binary-search based line indexing for large log files,
// enabling fast seeking to byte offsets by line number or timestamp prefix.
package lineindex

import (
	"bufio"
	"io"
	"sort"
)

// Entry records the byte offset and line number of a single log line.
type Entry struct {
	Line   int
	Offset int64
}

// Index is a sorted slice of Entries built from a ReadSeeker.
type Index []Entry

// Build scans r from the beginning and records the byte offset of every line.
// It resets r to the start before scanning.
func Build(r io.ReadSeeker) (Index, error) {
	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	var idx Index
	var offset int64
	lineNum := 1
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		idx = append(idx, Entry{Line: lineNum, Offset: offset})
		offset += int64(len(scanner.Bytes())) + 1 // +1 for newline
		lineNum++
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return idx, nil
}

// OffsetForLine returns the byte offset of the given 1-based line number.
// Returns -1 if the line number is out of range.
func (idx Index) OffsetForLine(line int) int64 {
	if line < 1 || line > len(idx) {
		return -1
	}
	return idx[line-1].Offset
}

// LineForOffset returns the 1-based line number whose offset is closest to
// (but not exceeding) the given byte offset using binary search.
func (idx Index) LineForOffset(offset int64) int {
	pos := sort.Search(len(idx), func(i int) bool {
		return idx[i].Offset > offset
	})
	if pos == 0 {
		return 1
	}
	return idx[pos-1].Line
}

// Len returns the total number of indexed lines.
func (idx Index) Len() int { return len(idx) }
