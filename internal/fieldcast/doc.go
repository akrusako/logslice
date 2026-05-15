// Package fieldcast provides a pipeline stage that coerces a named log field
// to a target scalar type (string, int, float, or bool).
//
// Both key=value and JSON log formats are supported. Lines where the target
// field is absent are forwarded unchanged. Lines where the conversion fails
// (e.g. the value is not numeric when casting to int) are also forwarded
// unchanged and the failure counter is incremented.
//
// Usage:
//
//	c, err := fieldcast.New("latency", fieldcast.TypeInt)
//	if err != nil { … }
//	outLine := c.Apply(inLine)
//	fmt.Println(c.Casted(), c.Failed())
package fieldcast
