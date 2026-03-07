//go:build darwin && arm64

package mclock

import "time"

var (
	freq        uint64
	useFallback bool
)

func init() {
	f := counterFreq()
	if f == 0 {
		useFallback = true
		return
	}
	freq = f
}

func counterValue() uint64 // implemented in clock_darwin_arm64.s
func counterFreq() uint64  // implemented in clock_darwin_arm64.s

// New creates a Clock anchored to the given epoch.
//
// The wall-clock elapsed time and the hardware counter are sampled
// sequentially, so a goroutine preemption between the two reads can
// introduce a small offset (typically nanoseconds, worst-case low
// microseconds under heavy GC pressure). This is well within the
// tolerance for millisecond timestamps and acceptable for microsecond
// timestamps in practice.
//
// If the hardware counter is unavailable, New falls back to [time.Since]
// using a re-anchored epoch with a monotonic reading.
func New(epoch time.Time) Clock {
	if useFallback {
		now := time.Now()
		return Clock{
			epoch: now.Add(epoch.Sub(now)),
		}
	}
	elapsed := time.Since(epoch)
	return Clock{
		baseTicks: counterValue(),
		baseMs:    elapsed.Milliseconds(),
		baseUs:    elapsed.Microseconds(),
	}
}

// Now returns milliseconds elapsed since the epoch.
func (c Clock) Now() int64 {
	if !c.epoch.IsZero() {
		return time.Since(c.epoch).Milliseconds()
	}
	delta := counterValue() - c.baseTicks
	return c.baseMs + ticksToMs(delta, freq)
}

// NowMicro returns microseconds elapsed since the epoch.
func (c Clock) NowMicro() int64 {
	if !c.epoch.IsZero() {
		return time.Since(c.epoch).Microseconds()
	}
	delta := counterValue() - c.baseTicks
	return c.baseUs + ticksToUs(delta, freq)
}
