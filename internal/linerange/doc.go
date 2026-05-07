// Package linerange provides parsing and evaluation of line number ranges
// for filtering log output by line position.
//
// A range string may take one of the following forms:
//
//	"N"    — exactly line N
//	"N:M"  — lines N through M (inclusive)
//	"N:"   — line N to end of file
//	":M"   — line 1 through M
//
// All line numbers are 1-based. Parsing rejects zero or negative values
// and ranges where End < Start.
package linerange
