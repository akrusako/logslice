// Package timeparse provides fast, format-agnostic timestamp parsing for
// logslice's log line processing pipeline.
//
// It supports the most common structured and semi-structured log timestamp
// formats (RFC3339, ISO 8601 variants, Apache/Nginx combined log format, and
// syslog-style timestamps) and caches the last successful format to minimise
// repeated format-trial overhead when processing large, homogeneously
// formatted log files.
//
// Typical usage:
//
//	p := timeparse.NewParser()
//	for _, line := range logLines {
//		ts, err := p.Parse(extractTimestamp(line))
//		if err != nil {
//			// handle unrecognised format
//		}
//		_ = ts
//	}
package timeparse
