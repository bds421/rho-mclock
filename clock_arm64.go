//go:build arm64

package mclock

import "time"

var (
	ticksPerMs uint64
	ticksPerUs uint64
)

func init() {
	freq := counterFreq()
	if freq == 0 {
		panic("mclock: CNTFRQ_EL0 returned 0")
	}
	ticksPerMs = freq / 1000
	if ticksPerMs == 0 {
		panic("mclock: counter frequency too low for millisecond precision")
	}
	ticksPerUs = freq / 1_000_000
	if ticksPerUs == 0 {
		panic("mclock: counter frequency too low for microsecond precision")
	}
}

func counterValue() uint64 // implemented in clock_arm64.s
func counterFreq() uint64  // implemented in clock_arm64.s

// New creates a Clock anchored to the given epoch.
//
// The wall-clock elapsed time and the hardware counter are sampled
// sequentially, so a goroutine preemption between the two reads can
// introduce a small offset (typically nanoseconds, worst-case low
// microseconds under heavy GC pressure). This is well within the
// tolerance for millisecond timestamps and acceptable for microsecond
// timestamps in practice.
func New(epoch time.Time) Clock {
	elapsed := time.Since(epoch)
	return Clock{
		baseTicks: counterValue(),
		baseMs:    elapsed.Milliseconds(),
		baseUs:    elapsed.Microseconds(),
	}
}

// Now returns milliseconds elapsed since the epoch.
func (c Clock) Now() int64 {
	// Unsigned subtraction is safe: CNTVCT_EL0 is a monotonically
	// increasing 64-bit counter that does not wrap during operation.
	delta := counterValue() - c.baseTicks
	return c.baseMs + int64(delta/ticksPerMs)
}

// NowMicro returns microseconds elapsed since the epoch.
func (c Clock) NowMicro() int64 {
	delta := counterValue() - c.baseTicks
	return c.baseUs + int64(delta/ticksPerUs)
}
