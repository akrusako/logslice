// Package timeparse provides utilities for parsing timestamps from log lines.
package timeparse

import (
	"fmt"
	"time"
)

// Common log timestamp formats ordered by specificity.
var knownFormats = []string{
	time.RFC3339Nano,
	time.RFC3339,
	"2006-01-02T15:04:05.999999999",
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05.999999999",
	"2006-01-02 15:04:05",
	"02/Jan/2006:15:04:05 -0700",
	"Jan 02 15:04:05",
}

// Parser holds a cached format for repeated parsing of similarly formatted logs.
type Parser struct {
	cachedFormat string
}

// NewParser returns a new Parser instance.
func NewParser() *Parser {
	return &Parser{}
}

// Parse attempts to extract a time.Time from a raw timestamp string.
// It caches the last successful format to speed up repeated calls.
func (p *Parser) Parse(raw string) (time.Time, error) {
	if p.cachedFormat != "" {
		if t, err := time.Parse(p.cachedFormat, raw); err == nil {
			return t, nil
		}
	}

	for _, format := range knownFormats {
		if t, err := time.Parse(format, raw); err == nil {
			p.cachedFormat = format
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("timeparse: unrecognized timestamp format: %q", raw)
}

// ParseAny attempts to parse a timestamp string without caching.
func ParseAny(raw string) (time.Time, error) {
	return NewParser().Parse(raw)
}

// Reset clears the cached format, forcing the next Parse call to try all known
// formats. This is useful when log lines from a different source with a
// different timestamp format are about to be parsed.
func (p *Parser) Reset() {
	p.cachedFormat = ""
}

// CachedFormat returns the format string that was last successfully used by
// Parse, or an empty string if no format has been cached yet.
func (p *Parser) CachedFormat() string {
	return p.cachedFormat
}
