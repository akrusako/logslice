// Package fieldpresence provides a line filter that accepts or rejects
// structured log lines based on the presence or absence of named fields.
//
// Lines are parsed as JSON objects or key=value pairs. A line passes when:
//   - every field listed in the "require" set is present, AND
//   - no field listed in the "forbid" set is present.
//
// Either constraint set may be empty, but at least one must be non-empty.
//
// Example:
//
//	c, err := fieldpresence.New([]string{"request_id"}, []string{"debug"})
//	if err != nil { ... }
//	if c.Allow(line) {
//	    fmt.Println(line)
//	}
package fieldpresence
