// Package timestamp provides utilities for extracting and formatting
// timestamps from structured log lines for display or comparison.
package timestamp

import (
	"strings"
	"time"

	"github.com/user/logslice/internal/fieldextract"
)

// Extractor extracts and optionally reformats timestamps from log lines.
type Extractor struct {
	fields    []string
	outFormat string
	strip     bool
}

// New creates an Extractor that looks for a timestamp in the given fields
// (in order) and reformats it using outFormat. If strip is true the original
// timestamp field is removed from the line before output.
func New(fields []string, outFormat string, strip bool) *Extractor {
	if len(fields) == 0 {
		fields = []string{"time", "ts", "timestamp", "@timestamp"}
	}
	if outFormat == "" {
		outFormat = time.RFC3339
	}
	return &Extractor{fields: fields, outFormat: outFormat, strip: strip}
}

// Extract returns the reformatted timestamp found in line, plus whether one
// was found. If strip is enabled the returned line has the timestamp field
// removed.
func (e *Extractor) Extract(line string) (ts string, out string, ok bool) {
	fields := fieldextract.Extract(line)
	for _, f := range e.fields {
		val, exists := fields[f]
		if !exists {
			continue
		}
		t, err := time.Parse(time.RFC3339Nano, val)
		if err != nil {
			t, err = time.Parse(time.RFC3339, val)
		}
		if err != nil {
			continue
		}
		ts = t.Format(e.outFormat)
		out = line
		if e.strip {
			out = stripField(line, f, val)
		}
		return ts, out, true
	}
	return "", line, false
}

// stripField removes key=value or "key":"value" occurrences from line.
func stripField(line, key, val string) string {
	// Try JSON-style removal first.
	jsonPat := `"` + key + `":"` + val + `"`
	if idx := strings.Index(line, jsonPat); idx >= 0 {
		rest := line[idx+len(jsonPat):]
		rest = strings.TrimPrefix(rest, ",")
		line = line[:idx] + rest
		return strings.TrimSpace(line)
	}
	// Try KV-style removal.
	kvPat := key + "=" + val
	if idx := strings.Index(line, kvPat); idx >= 0 {
		rest := line[idx+len(kvPat):]
		rest = strings.TrimPrefix(rest, " ")
		line = line[:idx] + rest
		return strings.TrimSpace(line)
	}
	return line
}
