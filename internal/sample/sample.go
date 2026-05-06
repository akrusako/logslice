// Package sample provides log line sampling strategies for reducing
// output volume while preserving statistical representativeness.
package sample

import (
	"math/rand"
	"sync"
)

// Sampler decides whether a given log line should be kept or dropped
// based on a configured sampling rate.
type Sampler struct {
	mu       sync.Mutex
	rate     float64 // 0.0 = drop all, 1.0 = keep all
	rng      *rand.Rand
	total    int64
	sampled  int64
}

// New creates a Sampler that keeps approximately rate*100 percent of lines.
// rate is clamped to [0.0, 1.0].
func New(rate float64, seed int64) *Sampler {
	if rate < 0 {
		rate = 0
	}
	if rate > 1 {
		rate = 1
	}
	return &Sampler{
		rate:    rate,
		rng:     rand.New(rand.NewSource(seed)),
	}
}

// Keep returns true if the line should be included in output.
func (s *Sampler) Keep(_ string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.total++
	if s.rate == 1.0 {
		s.sampled++
		return true
	}
	if s.rate == 0.0 {
		return false
	}
	if s.rng.Float64() < s.rate {
		s.sampled++
		return true
	}
	return false
}

// Stats returns total lines seen and lines kept.
func (s *Sampler) Stats() (total int64, kept int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.total, s.sampled
}

// Rate returns the configured sampling rate.
func (s *Sampler) Rate() float64 {
	return s.rate
}
