package mclock

import "time"

// Clock is a monotonic millisecond clock.
//
// On arm64, it reads the hardware counter register directly for ~5-8 ns
// timestamps. On other platforms, it falls back to [time.Since].
type Clock struct {
	baseTicks uint64    // arm64: CNTVCT_EL0 at creation
	baseMs    int64     // arm64: time.Since(epoch) at creation
	epoch     time.Time // fallback: stored epoch with monotonic reading
}
