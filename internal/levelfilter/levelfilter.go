// Package levelfilter provides log-level based filtering for structured log lines.
// It supports common log level names (debug, info, warn, error, fatal) and
// filters out lines whose level falls below the configured minimum severity.
package levelfilter

import "strings"

// Level represents a numeric log severity.
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
	LevelUnknown Level = -1
)

var levelNames = map[string]Level{
	"debug":   LevelDebug,
	"info":    LevelInfo,
	"warn":    LevelWarn,
	"warning": LevelWarn,
	"error":   LevelError,
	"err":     LevelError,
	"fatal":   LevelFatal,
	"crit":    LevelFatal,
	"critical": LevelFatal,
}

// Parse converts a level string to a Level value.
// Returns LevelUnknown if the string is not recognised.
func Parse(s string) Level {
	if l, ok := levelNames[strings.ToLower(strings.TrimSpace(s))]; ok {
		return l
	}
	return LevelUnknown
}

// Filter drops log lines whose severity is below the minimum level.
type Filter struct {
	min      Level
	field    string
	Dropped  int
	Passed   int
}

// New creates a Filter that passes lines at or above min severity.
// field is the key used to look up the level in key=value or JSON lines
// (e.g. "level" or "severity"). If field is empty, "level" is used.
func New(min Level, field string) *Filter {
	if field == "" {
		field = "level"
	}
	return &Filter{min: min, field: field}
}

// Allow returns true if the line should be kept.
// Lines that do not contain a recognisable level field are always passed through.
func (f *Filter) Allow(line string) bool {
	lvl := extractLevel(line, f.field)
	if lvl == LevelUnknown {
		f.Passed++
		return true
	}
	if lvl >= f.min {
		f.Passed++
		return true
	}
	f.Dropped++
	return false
}

// extractLevel scans line for field=value or "field":"value" patterns.
func extractLevel(line, field string) Level {
	// Try key=value: field=warn
	kv := field + "="
	if idx := strings.Index(line, kv); idx != -1 {
		rest := line[idx+len(kv):]
		val := strings.FieldsFunc(rest, func(r rune) bool {
			return r == ' ' || r == ',' || r == '"' || r == '}' || r == '\t'
		})
		if len(val) > 0 {
			return Parse(val[0])
		}
	}
	// Try JSON: "field":"warn"
	jsonKey := `"` + field + `":"`
	if idx := strings.Index(line, jsonKey); idx != -1 {
		rest := line[idx+len(jsonKey):]
		end := strings.IndexByte(rest, '"')
		if end != -1 {
			return Parse(rest[:end])
		}
	}
	return LevelUnknown
}
