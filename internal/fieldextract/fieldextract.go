// Package fieldextract parses structured log lines (JSON or key=value)
// and extracts named fields for filtering and formatting.
package fieldextract

import (
	"encoding/json"
	"strings"
)

// Fields is a map of field name to string value extracted from a log line.
type Fields map[string]string

// Extract attempts to parse a log line as JSON first, then falls back to
// key=value pair parsing. Returns an empty Fields map if neither succeeds.
func Extract(line string) Fields {
	line = strings.TrimSpace(line)
	if len(line) == 0 {
		return Fields{}
	}

	if line[0] == '{' {
		if f := extractJSON(line); len(f) > 0 {
			return f
		}
	}

	return extractKV(line)
}

// extractJSON parses a JSON object and converts top-level scalar values to strings.
func extractJSON(line string) Fields {
	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(line), &raw); err != nil {
		return Fields{}
	}

	f := make(Fields, len(raw))
	for k, v := range raw {
		switch val := v.(type) {
		case string:
			f[k] = val
		case nil:
			f[k] = ""
		default:
			b, err := json.Marshal(val)
			if err == nil {
				f[k] = string(b)
			}
		}
	}
	return f
}

// extractKV parses space-separated key=value tokens from a log line.
// Tokens that do not contain '=' are skipped.
func extractKV(line string) Fields {
	f := Fields{}
	for _, token := range strings.Fields(line) {
		idx := strings.IndexByte(token, '=')
		if idx <= 0 || idx == len(token)-1 {
			continue
		}
		key := token[:idx]
		val := strings.Trim(token[idx+1:], `"`)
		f[key] = val
	}
	return f
}
