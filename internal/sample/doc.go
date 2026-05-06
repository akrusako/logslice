// Package sample implements probabilistic log line sampling.
//
// A Sampler is created with a rate in [0.0, 1.0] and a random seed.
// Each call to Keep returns true with probability equal to the rate,
// allowing callers to reduce throughput while maintaining a statistically
// representative subset of log lines.
//
// Example:
//
//	s := sample.New(0.1, time.Now().UnixNano()) // keep ~10%
//	for scanner.Scan() {
//		if s.Keep(scanner.Text()) {
//			fmt.Println(scanner.Text())
//		}
//	}
package sample
