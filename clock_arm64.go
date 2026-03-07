//go:build arm64

package mclock

import "time"

var ticksPerMs uint64

func init() {
	freq := counterFreq()
	if freq == 0 {
		panic("mclock: CNTFRQ_EL0 returned 0")
	}
	ticksPerMs = freq / 1000
	if ticksPerMs == 0 {
		panic("mclock: counter frequency too low for millisecond precision")
	}
}

func counterValue() uint64 // implemented in clock_arm64.s
func counterFreq() uint64  // implemented in clock_arm64.s

// New creates a Clock anchored to the given epoch.
func New(epoch time.Time) Clock {
	return Clock{
		baseTicks: counterValue(),
		baseMs:    time.Since(epoch).Milliseconds(),
	}
}

// Now returns milliseconds elapsed since the epoch.
func (c Clock) Now() int64 {
	delta := counterValue() - c.baseTicks
	return c.baseMs + int64(delta/ticksPerMs)
}
