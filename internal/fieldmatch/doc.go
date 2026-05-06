// Package fieldmatch provides regex-based matching against individual fields
// extracted from structured log lines.
//
// It supports both JSON-encoded lines (e.g. {"level":"error","msg":"..."})
// and key=value logfmt-style lines (e.g. level=error msg="...").
//
// Usage:
//
//	m, err := fieldmatch.New(map[string]string{
//		"level":   "^(warn|error)$",
//		"service": "auth",
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	if m.Match(line) {
//		fmt.Println(line)
//	}
//
// All rules must be satisfied for a line to match (AND semantics).
// A Matcher with no rules matches every line.
package fieldmatch
