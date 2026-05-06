// Package transform provides line-level text transformation utilities for
// logslice pipelines.
//
// A Transformer is constructed with a set of Option values and then applied
// to each log line in sequence.  Supported transformations include:
//
//   - Upper / lower case folding
//   - Prefix and suffix stripping
//   - Field renaming for both JSON and key=value encoded lines
//
// Transformations are applied in the following order:
//  1. StripPrefix
//  2. StripSuffix
//  3. RenameFields
//  4. Case folding (upper takes precedence over lower)
//
// Example:
//
//	tr := transform.New(
//	    transform.WithStripPrefix("INFO "),
//	    transform.WithRenameFields(map[string]string{"msg": "message"}),
//	)
//	out := tr.Apply(line)
package transform
