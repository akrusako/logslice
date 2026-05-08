// Package jsonpath provides dot-notation field extraction and filtering
// for structured log lines (JSON and key=value formats).
package jsonpath

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Extractor resolves dot-notation paths against log lines.
type Extractor struct {
	path   []string
	raw    string
}

// New creates an Extractor for the given dot-notation path (e.g. "request.method").
func New(path string) (*Extractor, error) {
	if path == "" {
		return nil, fmt.Errorf("jsonpath: empty path")
	}
	parts := strings.Split(path, ".")
	for _, p := range parts {
		if p == "" {
			return nil, fmt.Errorf("jsonpath: invalid path %q (empty segment)", path)
		}
	}
	return &Extractor{path: parts, raw: path}, nil
}

// Get returns the string value at the configured path within line.
// Returns ("", false) when the path is not found or the line is not JSON.
func (e *Extractor) Get(line string) (string, bool) {
	line = strings.TrimSpace(line)
	if len(line) == 0 || line[0] != '{' {
		return "", false
	}
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return "", false
	}
	return traverse(obj, e.path)
}

// Path returns the original dot-notation path string.
func (e *Extractor) Path() string { return e.raw }

// traverse walks a nested map following the path segments.
func traverse(obj map[string]interface{}, path []string) (string, bool) {
	if len(path) == 0 {
		return "", false
	}
	v, ok := obj[path[0]]
	if !ok {
		return "", false
	}
	if len(path) == 1 {
		switch val := v.(type) {
		case string:
			return val, true
		case float64:
			return fmt.Sprintf("%g", val), true
		case bool:
			return fmt.Sprintf("%t", val), true
		default:
			b, err := json.Marshal(val)
			if err != nil {
				return "", false
			}
			return string(b), true
		}
	}
	nested, ok := v.(map[string]interface{})
	if !ok {
		return "", false
	}
	return traverse(nested, path[1:])
}
