// Package ratelimit implements a token-bucket rate limiter used by
// logslice to throttle log-line output.
//
// Two independent controls are provided:
//
//   - maxPerSec  – maximum lines emitted per second (token bucket).
//     Set to 0 to disable.
//   - maxTotal   – hard cap on the total number of lines emitted in
//     a single invocation. Set to 0 to disable.
//
// When both limits are active the stricter one wins: a line is only
// allowed when a rate token is available AND the total cap has not
// been reached.
package ratelimit
