package mclock

import "time"

// Clock is a monotonic millisecond clock.
//
// On darwin/arm64 (Apple Silicon), it reads the hardware counter register
// directly for ~5-8 ns timestamps. On other platforms, it falls back to
// [time.Since].
type Clock struct {
	baseTicks uint64    // arm64: CNTVCT_EL0 at creation
	baseMs    int64     // arm64: time.Since(epoch) at creation
	baseUs    int64     // arm64: time.Since(epoch) in microseconds at creation
	epoch     time.Time // fallback: stored epoch with monotonic reading
}

// ticksToMs converts a tick delta to milliseconds using exact
// quotient+remainder arithmetic. This avoids the drift that
// integer-truncated divisors (freq/1000) introduce on frequencies
// that are not exact multiples of 1000.
func ticksToMs(delta, freq uint64) int64 {
	return int64(delta/freq*1000 + delta%freq*1000/freq) // #nosec G115 — result is elapsed ms, overflow requires ~292M years
}

// ticksToUs converts a tick delta to microseconds using exact
// quotient+remainder arithmetic.
func ticksToUs(delta, freq uint64) int64 {
	return int64(delta/freq*1_000_000 + delta%freq*1_000_000/freq) // #nosec G115 — result is elapsed µs, overflow requires ~292K years
}
