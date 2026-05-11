// Package split provides a line splitter that emits each log line into one
// of N buckets based on a modulo hash of the line content. This is useful for
// sharding large log files across parallel workers.
package split

import (
	"fmt"
	"hash/fnv"
)

// Splitter assigns lines to buckets.
type Splitter struct {
	buckets  uint64
	bucket   uint64
	total    int
	matched  int
}

// New creates a Splitter that keeps lines whose hash falls in the given bucket
// (0-indexed) out of the total number of buckets.
// Returns an error if buckets < 1 or bucket >= buckets.
func New(buckets, bucket uint64) (*Splitter, error) {
	if buckets < 1 {
		return nil, fmt.Errorf("split: buckets must be >= 1, got %d", buckets)
	}
	if bucket >= buckets {
		return nil, fmt.Errorf("split: bucket index %d out of range [0, %d)", bucket, buckets)
	}
	return &Splitter{buckets: buckets, bucket: bucket}, nil
}

// Allow returns true if the line belongs to this splitter's bucket.
func (s *Splitter) Allow(line string) bool {
	s.total++
	h := hash(line)
	if h%s.buckets == s.bucket {
		s.matched++
		return true
	}
	return false
}

// Total returns the number of lines evaluated.
func (s *Splitter) Total() int { return s.total }

// Matched returns the number of lines assigned to this bucket.
func (s *Splitter) Matched() int { return s.matched }

// Bucket returns the zero-based bucket index this splitter represents.
func (s *Splitter) Bucket() uint64 { return s.bucket }

// Buckets returns the total number of buckets.
func (s *Splitter) Buckets() uint64 { return s.buckets }

func hash(line string) uint64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(line))
	return h.Sum64()
}
