// Package highlight provides lightweight ANSI terminal highlighting for
// logslice output.
//
// A Highlighter can be constructed with a specific colour code and an enabled
// flag. When disabled every method is a no-op, so callers do not need to
// guard calls behind their own conditional logic.
//
// Typical usage:
//
//	h := highlight.New(highlight.Yellow, isTTY)
//	fmt.Println(h.Terms(line, searchTerms))
//
// Predefined colour constants (Red, Green, Yellow, Cyan, Bold, Reset) cover
// the most common terminal palette without pulling in external dependencies.
package highlight
