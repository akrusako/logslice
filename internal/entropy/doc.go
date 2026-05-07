// Package entropy provides Shannon-entropy scoring for log lines.
//
// Typical usage:
//
//	sc := entropy.New(4.5) // drop lines with entropy >= 4.5 bits/byte
//	for _, line := range lines {
//	    if sc.Allow(line) {
//	        fmt.Println(line)
//	    }
//	}
//	log.Printf("dropped %d/%d high-entropy lines", sc.Above(), sc.Total())
//
// A threshold of 0 disables filtering; every line is passed through.
// The Score function is also exported for use in pipelines or reporting.
package entropy
