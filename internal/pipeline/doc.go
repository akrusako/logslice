// Package pipeline provides a composable, ordered line-processing chain.
//
// Each Stage in the pipeline receives the current line text and decides
// whether to pass it downstream (optionally transformed) or discard it.
// Stages are applied left-to-right; the first stage that drops a line
// short-circuits the remaining stages.
//
// Typical usage:
//
//	p := pipeline.New(os.Stdout,
//	    filterStage,
//	    transformStage,
//	    highlightStage,
//	)
//	for _, line := range lines {
//	    p.Process(line)
//	}
//	passed, dropped := p.Stats()
package pipeline
