// Package redact replaces sensitive substrings in log lines using compiled
// regular expressions.
//
// # Usage
//
//	r, err := redact.New([]string{`\d{4}-\d{4}-\d{4}-\d{4}`}, "[CARD]")
//	if err != nil {
//		log.Fatal(err)
//	}
//	clean, changed := r.Apply(line)
//
// Multiple patterns can be registered; all are applied in order. The mask
// string defaults to "[REDACTED]" when left empty. Count() reports how many
// lines were modified.
package redact
