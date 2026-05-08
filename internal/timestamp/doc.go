// Package timestamp provides structured log timestamp extraction and
// reformatting.
//
// It inspects log lines (JSON or key=value format) for common timestamp
// fields such as "time", "ts", "timestamp", and "@timestamp", parses the
// value as RFC3339 or RFC3339Nano, and re-emits it in any desired Go
// time-format layout.
//
// Optionally the original timestamp field can be stripped from the output
// line so downstream consumers receive a cleaner record.
//
// Example:
//
//	e := timestamp.New(nil, time.RFC3339, false)
//	ts, line, ok := e.Extract(`{"time":"2024-01-01T00:00:00Z","msg":"hi"}`)
//	// ts  == "2024-01-01T00:00:00Z"
//	// ok  == true
package timestamp
