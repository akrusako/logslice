// Package labelinject provides an Injector that appends static key=value
// labels to every processed log line.
//
// Labels are injected in a format-aware manner:
//   - JSON objects receive the labels as additional top-level fields.
//   - Key=value lines receive the labels appended as "key=value" pairs.
//   - Plain-text lines receive a space-separated "key=value" suffix.
//
// Typical usage:
//
//	inj, err := labelinject.New([]string{"env=prod", "region=us-east-1"})
//	if err != nil { ... }
//	output := inj.Apply(inputLine)
//
// The Injected() counter reports how many lines were modified.
package labelinject
