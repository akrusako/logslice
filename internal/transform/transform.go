// Package transform applies line-level text transformations such as
// case folding, prefix/suffix stripping, and field renaming.
package transform

import (
	"encoding/json"
	"strings"
)

// Transformer holds configuration for a set of transformations.
type Transformer struct {
	upper      bool
	lower      bool
	stripPrefix string
	stripSuffix string
	renameFields map[string]string
}

// Option configures a Transformer.
type Option func(*Transformer)

// WithUpper converts each line to upper case.
func WithUpper() Option { return func(t *Transformer) { t.upper = true } }

// WithLower converts each line to lower case.
func WithLower() Option { return func(t *Transformer) { t.lower = true } }

// WithStripPrefix removes a fixed prefix from each line when present.
func WithStripPrefix(p string) Option { return func(t *Transformer) { t.stripPrefix = p } }

// WithStripSuffix removes a fixed suffix from each line when present.
func WithStripSuffix(s string) Option { return func(t *Transformer) { t.stripSuffix = s } }

// WithRenameFields renames JSON or key=value fields according to the map.
func WithRenameFields(m map[string]string) Option {
	return func(t *Transformer) { t.renameFields = m }
}

// New creates a Transformer with the supplied options.
func New(opts ...Option) *Transformer {
	t := &Transformer{}
	for _, o := range opts {
		o(t)
	}
	return t
}

// Apply runs all configured transformations on line and returns the result.
func (t *Transformer) Apply(line string) string {
	if t.stripPrefix != "" {
		line = strings.TrimPrefix(line, t.stripPrefix)
	}
	if t.stripSuffix != "" {
		line = strings.TrimSuffix(line, t.stripSuffix)
	}
	if len(t.renameFields) > 0 {
		line = t.applyRename(line)
	}
	if t.upper {
		line = strings.ToUpper(line)
	} else if t.lower {
		line = strings.ToLower(line)
	}
	return line
}

func (t *Transformer) applyRename(line string) string {
	// Try JSON first.
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err == nil {
		for old, new := range t.renameFields {
			if v, ok := obj[old]; ok {
				delete(obj, old)
				obj[new] = v
			}
		}
		if b, err := json.Marshal(obj); err == nil {
			return string(b)
		}
	}
	// Fall back to key=value.
	parts := strings.Fields(line)
	for i, p := range parts {
		if idx := strings.IndexByte(p, '='); idx > 0 {
			key := p[:idx]
			if newKey, ok := t.renameFields[key]; ok {
				parts[i] = newKey + p[idx:]
			}
		}
	}
	return strings.Join(parts, " ")
}
