// Package rotate detects and handles log file rotation, allowing
// continuous reading across file truncation or replacement events.
package rotate

import (
	"errors"
	"io"
	"os"
	"time"
)

// Watcher monitors a log file for rotation events (truncation or inode
// change) and re-opens the file when rotation is detected.
type Watcher struct {
	path     string
	file     *os.File
	lastSize int64
	lastIno  uint64
	pollInterval time.Duration
}

// New opens the named file and returns a Watcher ready to detect rotation.
// pollInterval controls how often the file is checked for rotation.
func New(path string, pollInterval time.Duration) (*Watcher, error) {
	if pollInterval <= 0 {
		pollInterval = 500 * time.Millisecond
	}
	w := &Watcher{path: path, pollInterval: pollInterval}
	if err := w.open(); err != nil {
		return nil, err
	}
	return w, nil
}

// open (re-)opens the file and records its current size and inode.
func (w *Watcher) open() error {
	f, err := os.Open(w.path)
	if err != nil {
		return err
	}
	if w.file != nil {
		_ = w.file.Close()
	}
	w.file = f
	info, err := f.Stat()
	if err != nil {
		return err
	}
	w.lastSize = info.Size()
	w.lastIno = inode(info)
	return nil
}

// Rotated returns true if the file has been rotated since the last check.
// It updates internal state as a side-effect.
func (w *Watcher) Rotated() (bool, error) {
	info, err := os.Stat(w.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return true, nil
		}
		return false, err
	}
	curIno := inode(info)
	curSize := info.Size()
	if curIno != w.lastIno || curSize < w.lastSize {
		w.lastIno = curIno
		w.lastSize = curSize
		return true, nil
	}
	w.lastSize = curSize
	return false, nil
}

// Reopen closes the current file handle and opens a fresh one.
func (w *Watcher) Reopen() error {
	return w.open()
}

// Reader returns an io.Reader positioned at the current file offset.
func (w *Watcher) Reader() io.Reader {
	return w.file
}

// Close releases the underlying file handle.
func (w *Watcher) Close() error {
	if w.file != nil {
		return w.file.Close()
	}
	return nil
}

// PollInterval returns the configured poll interval.
func (w *Watcher) PollInterval() time.Duration {
	return w.pollInterval
}
