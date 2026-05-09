// Package prefix provides a lightweight Prefixer that prepends a configurable
// label to every non-empty log line passing through a pipeline stage.
//
// Typical usage:
//
//	p := prefix.New("[app] ")
//	for _, line := range lines {
//		fmt.Println(p.Apply(line))
//	}
//
// The prefix can be changed at runtime with SetPrefix, and Strip can be used
// to reverse the operation when reading prefixed logs.
package prefix
