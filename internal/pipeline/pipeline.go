// Package pipeline wires together a sequence of line-processing stages
// into a single reusable processing chain.
package pipeline

import "io"

// Stage is a function that receives a line and returns the (possibly
// transformed) line and whether it should be kept in the output stream.
type Stage func(line string) (out string, keep bool)

// Pipeline applies an ordered list of Stages to every line written to it
// and forwards surviving lines to the underlying writer.
type Pipeline struct {
	w      io.Writer
	stages []Stage
	passed int
	dropped int
}

// New creates a Pipeline that writes accepted lines to w.
func New(w io.Writer, stages ...Stage) *Pipeline {
	return &Pipeline{w: w, stages: stages}
}

// Process runs line through every stage in order. If any stage drops the
// line the function returns false without writing. Otherwise the final
// transformed line (with a trailing newline) is written to the underlying
// writer and the function returns true.
func (p *Pipeline) Process(line string) (bool, error) {
	current := line
	for _, s := range p.stages {
		out, keep := s(current)
		if !keep {
			p.dropped++
			return false, nil
		}
		current = out
	}
	_, err := io.WriteString(p.w, current+"\n")
	if err != nil {
		return false, err
	}
	p.passed++
	return true, nil
}

// Stats returns the number of lines that passed and were dropped.
func (p *Pipeline) Stats() (passed, dropped int) {
	return p.passed, p.dropped
}
