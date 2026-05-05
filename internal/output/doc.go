// Package output provides log line formatting for logslice output.
//
// It supports three output formats:
//
//   - FormatRaw:  Lines are written as-is (default). Suitable for piping
//     to other tools or preserving the original log format.
//
//   - FormatJSON: JSON log lines are pretty-printed with indentation.
//     Non-JSON lines are passed through unchanged.
//
//   - FormatTSV:  JSON log lines are rendered as tab-separated key=value
//     pairs, useful for import into spreadsheets or awk processing.
//     Non-JSON lines are passed through unchanged.
//
// Example usage:
//
//	f := output.New(os.Stdout, output.FormatJSON)
//	f.WriteLine(`{"level":"info","msg":"started"}`)
package output
