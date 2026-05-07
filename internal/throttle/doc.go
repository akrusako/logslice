// Package throttle implements a sliding-window line-rate throttle for log
// stream processing.
//
// A Throttle is constructed with a window duration and a maximum number of
// lines permitted within that window. Each call to Allow advances an internal
// counter; once the cap is reached every subsequent call returns false until
// old entries slide out of the window.
//
// Throttle is safe for concurrent use.
//
// Example:
//
//	th := throttle.New(time.Second, 1000)
//	for _, line := range lines {
//		if th.Allow() {
//			fmt.Println(line)
//		}
//	}
package throttle
