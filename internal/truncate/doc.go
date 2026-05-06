// Package truncate implements configurable line-length truncation for
// log processing pipelines.
//
// A Truncator is constructed with a maximum byte length. Lines that
// exceed the limit are trimmed and a "..." suffix is appended so the
// result fits within the configured size.
//
// When maxLen is zero or negative the Truncator is disabled and every
// line passes through unchanged, making it safe to wire unconditionally
// into a processing chain.
//
// Usage:
//
//	tr := truncate.New(120)
//	for _, line := range lines {
//		fmt.Println(tr.Apply(line))
//	}
//	fmt.Println(tr.Summary())
package truncate
