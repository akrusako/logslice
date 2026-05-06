// Package linecount provides fast line counting for log files,
// supporting both full-file counts and counted reads from an offset.
package linecount

import (
	"bufio"
	"io"
	"os"
)

// Result holds the outcome of a line count operation.
type Result struct {
	Lines int64
	Bytes int64
}

// File counts all newline-terminated lines in the named file.
// It returns the total line count and total byte size, or an error.
func File(path string) (Result, error) {
	f, err := os.Open(path)
	if err != nil {
		return Result{}, err
	}
	defer f.Close()
	return Reader(f)
}

// Reader counts lines from any io.Reader.
// Each newline character ('\n') increments the counter.
// A final line without a trailing newline is also counted.
func Reader(r io.Reader) (Result, error) {
	var lines, bytes int64
	br := bufio.NewReaderSize(r, 64*1024)
	var lastByte byte
	for {
		b, err := br.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			return Result{}, err
		}
		bytes++
		if b == '\n' {
			lines++
		}
		lastByte = b
	}
	// Count a non-empty final line that lacks a trailing newline.
	if bytes > 0 && lastByte != '\n' {
		lines++
	}
	return Result{Lines: lines, Bytes: bytes}, nil
}

// FromOffset counts lines in the file starting at the given byte offset.
// Useful for incremental counting after a previous read position.
func FromOffset(path string, offset int64) (Result, error) {
	f, err := os.Open(path)
	if err != nil {
		return Result{}, err
	}
	defer f.Close()
	if _, err := f.Seek(offset, io.SeekStart); err != nil {
		return Result{}, err
	}
	return Reader(f)
}
