// Package columnar extracts named fields from structured log lines and
// renders them as aligned, delimited columns.
//
// It supports both JSON and key=value log formats via the fieldextract
// package. Missing fields are represented as empty strings and counted
// separately so callers can decide whether to skip incomplete rows.
//
// Example usage:
//
//	c := columnar.New([]string{"level", "msg", "ts"}, 12, "\t")
//	fmt.Println(c.Header())
//	for _, line := range lines {
//		row, _ := c.Format(line)
//		fmt.Println(row)
//	}
package columnar
