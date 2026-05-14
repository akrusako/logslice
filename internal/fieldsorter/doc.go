// Package fieldsorter reorders fields within structured log lines.
//
// It accepts a list of field names that should appear first in the
// output, in the specified order. Any fields not listed are appended
// after the ordered fields in their original relative order.
//
// Both JSON object lines and key=value pair lines are supported.
// Lines that cannot be parsed as either format are passed through
// unchanged and counted as skipped.
//
// Example usage:
//
//	s, err := fieldsorter.New([]string{"level", "msg", "ts"})
//	if err != nil {
//		log.Fatal(err)
//	}
//	out := s.Apply(`{"ts":"…","msg":"hello","level":"info"}`)
//	// out: {"level":"info","msg":"hello","ts":"…"}
package fieldsorter
