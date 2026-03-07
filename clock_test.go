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

func TestNowMicroPositive(t *testing.T) {
	epoch := time.Now().Add(-time.Second)
	clk := New(epoch)
	us := clk.NowMicro()
	if us <= 0 {
		t.Fatalf("expected positive us for past epoch, got %d", us)
	}
}

func TestNowMicroMonotonic(t *testing.T) {
	clk := New(time.Now().Add(-time.Hour))
	prev := clk.NowMicro()
	for i := range 10_000 {
		cur := clk.NowMicro()
		if cur < prev {
			t.Fatalf("non-monotonic at iteration %d: %d < %d", i, cur, prev)
		}
		prev = cur
	}
}

func TestNowMicroAccuracy(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping accuracy test in short mode")
	}
	epoch := time.Now().Add(-time.Hour)
	clk := New(epoch)
	before := clk.NowMicro()
	time.Sleep(10 * time.Millisecond)
	after := clk.NowMicro()
	delta := after - before
	if delta < 8000 || delta > 12000 {
		t.Fatalf("expected ~10000us delta, got %dus", delta)
	}
}

func TestFutureEpoch(t *testing.T) {
	epoch := time.Now().Add(time.Second)
	clk := New(epoch)
	ms := clk.Now()
	if ms >= 0 {
		t.Fatalf("expected negative ms for future epoch, got %d", ms)
	}
}

func TestNowMicroDivByThousandApproxNow(t *testing.T) {
	epoch := time.Now().Add(-time.Hour)
	clk := New(epoch)
	for range 100 {
		ms := clk.Now()
		us := clk.NowMicro()
		usDiv := us / 1000
		diff := usDiv - ms
		if diff < -1 || diff > 1 {
			t.Fatalf("NowMicro()/1000=%d, Now()=%d, diff=%d (want within +/-1)", usDiv, ms, diff)
		}
	}
}

func TestTicksToMs(t *testing.T) {
	tests := []struct {
		name        string
		delta, freq uint64
		want        int64
	}{
		{"zero delta", 0, 24_000_000, 0},
		{"24MHz 1ms", 24_000, 24_000_000, 1},
		{"24MHz 1s", 24_000_000, 24_000_000, 1000},
		{"24MHz 1000s", 24_000_000_000, 24_000_000, 1_000_000},
		{"19.2MHz 1ms", 19_200, 19_200_000, 1},
		{"19.2MHz 1s", 19_200_000, 19_200_000, 1000},
		{"19.2MHz 500ms", 9_600_000, 19_200_000, 500},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ticksToMs(tt.delta, tt.freq)
			if got != tt.want {
				t.Errorf("ticksToMs(%d, %d) = %d, want %d", tt.delta, tt.freq, got, tt.want)
			}
		})
	}
}

func TestTicksToUs(t *testing.T) {
	tests := []struct {
		name        string
		delta, freq uint64
		want        int64
	}{
		{"zero delta", 0, 24_000_000, 0},
		{"24MHz 1us", 24, 24_000_000, 1},
		{"24MHz 1s", 24_000_000, 24_000_000, 1_000_000},
		{"24MHz 1000s", 24_000_000_000, 24_000_000, 1_000_000_000},
		// 19.2MHz: freq/1_000_000 = 19 (truncated), old code gives
		// 19200/19 = 1010us. Exact arithmetic gives 1000us.
		{"19.2MHz 1ms exact", 19_200, 19_200_000, 1000},
		{"19.2MHz 1s", 19_200_000, 19_200_000, 1_000_000},
		{"19.2MHz 500ms", 9_600_000, 19_200_000, 500_000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ticksToUs(tt.delta, tt.freq)
			if got != tt.want {
				t.Errorf("ticksToUs(%d, %d) = %d, want %d", tt.delta, tt.freq, got, tt.want)
			}
		})
	}
}

func TestLongRunAccuracy(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long-run accuracy test in short mode")
	}
	epoch := time.Now().Add(-time.Hour)
	clk := New(epoch)
	start := time.Now()
	startMs := clk.Now()
	time.Sleep(5 * time.Second)
	elapsed := time.Since(start).Milliseconds()
	delta := clk.Now() - startMs
	drift := delta - elapsed
	if drift < -2 || drift > 2 {
		t.Fatalf("drift after 5s: mclock=%dms, time.Since=%dms, drift=%dms", delta, elapsed, drift)
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

func BenchmarkNowMicro(b *testing.B) {
	clk := New(time.Now().Add(-time.Hour))
	b.ReportAllocs()
	for b.Loop() {
		sink = clk.NowMicro()
	}
}

func BenchmarkTimeSince(b *testing.B) {
	epoch := time.Now().Add(-time.Hour)
	b.ReportAllocs()
	for b.Loop() {
		sink = time.Since(epoch).Milliseconds()
	}
}

func BenchmarkTimeSinceMicro(b *testing.B) {
	epoch := time.Now().Add(-time.Hour)
	b.ReportAllocs()
	for b.Loop() {
		sink = time.Since(epoch).Microseconds()
	}
}
