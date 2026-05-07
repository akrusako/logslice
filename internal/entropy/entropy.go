// Package entropy measures the Shannon entropy of log lines to detect
// anomalous or high-information-density entries (e.g. encoded payloads,
// stack traces, or binary data smuggled into logs).
package entropy

import (
	"math"
)

// Scorer computes per-line entropy and emits lines that exceed a threshold.
type Scorer struct {
	threshold float64
	above     int
	total     int
}

// New returns a Scorer that flags lines whose Shannon entropy exceeds
// threshold bits-per-byte. A threshold of 0 disables filtering (all lines
// pass through).
func New(threshold float64) *Scorer {
	return &Scorer{threshold: threshold}
}

// Score returns the Shannon entropy of s in bits per byte.
func Score(s string) float64 {
	if len(s) == 0 {
		return 0
	}
	var freq [256]int
	for i := 0; i < len(s); i++ {
		freq[s[i]]++
	}
	n := float64(len(s))
	var h float64
	for _, c := range freq {
		if c == 0 {
			continue
		}
		p := float64(c) / n
		h -= p * math.Log2(p)
	}
	return h
}

// Allow returns true when the line should pass through.
// If the threshold is zero every line is allowed. Otherwise only lines
// whose entropy is strictly below the threshold are allowed.
func (s *Scorer) Allow(line string) bool {
	s.total++
	if s.threshold == 0 {
		return true
	}
	if Score(line) >= s.threshold {
		s.above++
		return false
	}
	return true
}

// Above returns the number of lines dropped for exceeding the threshold.
func (s *Scorer) Above() int { return s.above }

// Total returns the total number of lines evaluated.
func (s *Scorer) Total() int { return s.total }
