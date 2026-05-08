// Package jsonpath provides dot-notation path extraction for structured
// log lines stored as JSON objects.
//
// Paths use period-separated segments to navigate nested objects:
//
//	request.method  →  {"request":{"method":"GET"}}
//	status          →  {"status":200}
//
// Only JSON object lines are supported; key=value lines and plain text
// lines will always return a not-found result.
//
// Usage:
//
//	e, err := jsonpath.New("request.method")
//	if err != nil { ... }
//	val, ok := e.Get(line)
package jsonpath
