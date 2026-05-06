// Package dedupe provides line-level deduplication for log streams.
// It tracks recently seen lines using a fixed-size ring buffer of hashes
// and suppresses repeated entries within the configured window.
package dedupe

import (
	"hash/fnv"
)

// Deduper suppresses duplicate log lines within a sliding window.
type Deduper struct {
	window []uint64
	size   int
	pos    int
	seen   map[uint64]struct{}
}

// New creates a Deduper with the given window size.
// A window of 0 disables deduplication (all lines pass through).
func New(windowSize int) *Deduper {
	if windowSize <= 0 {
		return &Deduper{}
	}
	return &Deduper{
		window: make([]uint64, windowSize),
		size:   windowSize,
		seen:   make(map[uint64]struct{}, windowSize),
	}
}

// IsDuplicate returns true if line was seen within the current window.
// If it is not a duplicate, the line is recorded and true is returned on
// future calls within the window.
func (d *Deduper) IsDuplicate(line string) bool {
	if d.size == 0 {
		return false
	}
	h := hash(line)
	if _, ok := d.seen[h]; ok {
		return true
	}
	// Evict the oldest entry.
	old := d.window[d.pos]
	delete(d.seen, old)
	// Insert new.
	d.window[d.pos] = h
	d.seen[h] = struct{}{}
	d.pos = (d.pos + 1) % d.size
	return false
}

// Reset clears all tracked state.
func (d *Deduper) Reset() {
	for k := range d.seen {
		delete(d.seen, k)
	}
	for i := range d.window {
		d.window[i] = 0
	}
	d.pos = 0
}

func hash(s string) uint64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(s))
	return h.Sum64()
}
