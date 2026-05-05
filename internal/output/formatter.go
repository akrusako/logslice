// Package output provides formatting utilities for log output.
package output

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// Format represents the output format for log lines.
type Format int

const (
	FormatRaw  Format = iota // Output lines as-is
	FormatJSON               // Pretty-print JSON lines
	FormatTSV                // Tab-separated key=value pairs
)

// Formatter writes formatted log lines to an output writer.
type Formatter struct {
	w      io.Writer
	format Format
}

// New creates a new Formatter writing to w in the given format.
func New(w io.Writer, format Format) *Formatter {
	return &Formatter{w: w, format: format}
}

// WriteLine writes a single log line according to the configured format.
func (f *Formatter) WriteLine(line string) error {
	switch f.format {
	case FormatJSON:
		return f.writeJSON(line)
	case FormatTSV:
		return f.writeTSV(line)
	default:
		_, err := fmt.Fprintln(f.w, line)
		return err
	}
}

func (f *Formatter) writeJSON(line string) error {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		// Not valid JSON — fall back to raw output
		_, err = fmt.Fprintln(f.w, line)
		return err
	}
	out, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(f.w, string(out))
	return err
}

func (f *Formatter) writeTSV(line string) error {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		_, err = fmt.Fprintln(f.w, line)
		return err
	}
	parts := make([]string, 0, len(obj))
	for k, v := range obj {
		parts = append(parts, fmt.Sprintf("%s=%v", k, v))
	}
	_, err := fmt.Fprintln(f.w, strings.Join(parts, "\t"))
	return err
}
