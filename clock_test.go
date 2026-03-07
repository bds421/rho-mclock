package mclock

import (
	"sync"
	"testing"
	"time"
)

func TestNowPositive(t *testing.T) {
	epoch := time.Now().Add(-time.Second)
	clk := New(epoch)
	ms := clk.Now()
	if ms <= 0 {
		t.Fatalf("expected positive ms for past epoch, got %d", ms)
	}
}

func TestNowMonotonic(t *testing.T) {
	clk := New(time.Now().Add(-time.Hour))
	prev := clk.Now()
	for i := range 10_000 {
		cur := clk.Now()
		if cur < prev {
			t.Fatalf("non-monotonic at iteration %d: %d < %d", i, cur, prev)
		}
		prev = cur
	}
}

func TestNowAccuracy(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping accuracy test in short mode")
	}
	epoch := time.Now().Add(-time.Hour)
	clk := New(epoch)
	before := clk.Now()
	time.Sleep(100 * time.Millisecond)
	after := clk.Now()
	delta := after - before
	if delta < 80 || delta > 120 {
		t.Fatalf("expected ~100ms delta, got %dms", delta)
	}
}

func TestNowConcurrent(t *testing.T) {
	clk := New(time.Now().Add(-time.Hour))
	var wg sync.WaitGroup
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 10_000 {
				ms := clk.Now()
				if ms <= 0 {
					t.Errorf("expected positive ms, got %d", ms)
					return
				}
			}
		}()
	}
	wg.Wait()
}

func TestNowAgreesWithTimeSince(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping agreement test in short mode")
	}
	epoch := time.Now().Add(-time.Hour)
	clk := New(epoch)

	// Compare over a short window — results should agree within +/- 1ms.
	for range 100 {
		mclockMs := clk.Now()
		stdMs := time.Since(epoch).Milliseconds()
		diff := mclockMs - stdMs
		if diff < -1 || diff > 1 {
			t.Fatalf("mclock=%d, time.Since=%d, diff=%d (want within +/-1)", mclockMs, stdMs, diff)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

var sink int64

func BenchmarkNow(b *testing.B) {
	clk := New(time.Now().Add(-time.Hour))
	b.ReportAllocs()
	for b.Loop() {
		sink = clk.Now()
	}
}

func BenchmarkTimeSince(b *testing.B) {
	now := time.Now()
	epoch := now.Add(time.Now().Add(-time.Hour).Sub(now))
	b.ReportAllocs()
	for b.Loop() {
		sink = time.Since(epoch).Milliseconds()
	}
}
