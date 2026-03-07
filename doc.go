// Package mclock provides a fast monotonic clock with millisecond and
// microsecond resolution.
//
// On arm64 platforms, mclock reads the CNTVCT_EL0 counter register directly
// via a single MRS instruction (~5-8 ns), bypassing Go's runtime and libSystem.
// On other architectures, it falls back to [time.Since].
//
// Usage:
//
//	clk := mclock.New(epoch)
//	ms := clk.Now()      // milliseconds since epoch
//	us := clk.NowMicro() // microseconds since epoch
package mclock
