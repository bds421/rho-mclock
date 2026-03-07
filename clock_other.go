//go:build !arm64

package mclock

import "time"

// New creates a Clock anchored to the given epoch.
func New(epoch time.Time) Clock {
	now := time.Now()
	return Clock{
		epoch: now.Add(epoch.Sub(now)),
	}
}

// Now returns milliseconds elapsed since the epoch.
func (c Clock) Now() int64 {
	return time.Since(c.epoch).Milliseconds()
}
