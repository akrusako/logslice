// Package context defines PipelineCtx, a thin wrapper around the standard
// library context.Context that carries logslice-specific configuration and
// per-run counters through each stage of a processing pipeline.
//
// Usage:
//
//	cfg := context.Config{MaxLines: 500, Quiet: false}
//	pc  := context.New(ctx, cfg)
//
//	for _, line := range lines {
//		if !pc.RecordLine() {
//			break // MaxLines reached
//		}
//		// … process line …
//	}
//
// PipelineCtx embeds context.Context so it can be passed wherever a standard
// context is expected, enabling cooperative cancellation across goroutines.
package context
