// Package context provides a pipeline context that carries cancellation,
// shared configuration, and per-run metadata through logslice processing stages.
package context

import (
	"context"
	"time"
)

// Config holds the shared runtime configuration for a logslice pipeline run.
type Config struct {
	// MaxLines caps total lines emitted; 0 means unlimited.
	MaxLines int
	// Since and Until define the optional time-range filter window.
	Since time.Time
	Until time.Time
	// Quiet suppresses progress/stats output.
	Quiet bool
}

// PipelineCtx wraps a standard context.Context with logslice-specific values.
type PipelineCtx struct {
	context.Context
	Cfg    Config
	Start  time.Time
	linesOut int
}

// New creates a PipelineCtx derived from parent with the given Config.
func New(parent context.Context, cfg Config) *PipelineCtx {
	return &PipelineCtx{
		Context: parent,
		Cfg:     cfg,
		Start:   time.Now(),
	}
}

// RecordLine increments the emitted-line counter and cancels the context when
// MaxLines is reached, returning false to signal the caller to stop.
func (p *PipelineCtx) RecordLine() bool {
	p.linesOut++
	if p.Cfg.MaxLines > 0 && p.linesOut >= p.Cfg.MaxLines {
		return false
	}
	return true
}

// LinesOut returns the number of lines recorded so far.
func (p *PipelineCtx) LinesOut() int { return p.linesOut }

// Elapsed returns the duration since the pipeline started.
func (p *PipelineCtx) Elapsed() time.Duration { return time.Since(p.Start) }

// InTimeRange returns true when t falls within [Since, Until].
// If a bound is zero it is treated as open.
func (p *PipelineCtx) InTimeRange(t time.Time) bool {
	if !p.Cfg.Since.IsZero() && t.Before(p.Cfg.Since) {
		return false
	}
	if !p.Cfg.Until.IsZero() && t.After(p.Cfg.Until) {
		return false
	}
	return true
}
