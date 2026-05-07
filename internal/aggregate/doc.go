// Package aggregate counts and groups log lines by a named field value.
//
// It is useful for producing quick frequency summaries such as how many
// lines were emitted at each log level, or how many requests hit each
// HTTP status code.
//
// Usage:
//
//	c := aggregate.New("level")
//	for _, line := range lines {
//		c.Record(line)
//	}
//	c.WriteSummary(os.Stdout)
//
// The field value is extracted via fieldextract.Extract, so both
// key=value and JSON formats are supported automatically.
package aggregate
