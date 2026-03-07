# mclock

Fast monotonic millisecond clock for Go.

On arm64, reads the `CNTVCT_EL0` counter register directly via a single `MRS` instruction (~1.8 ns), bypassing Go's runtime and libSystem. Falls back to `time.Since` on other architectures.

## Usage

```go
import "github.com/bds421/rho-mclock"

epoch := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
clk := mclock.New(epoch)
ms := clk.Now() // milliseconds since epoch
```

## How It Works

**arm64 fast path:**
1. `New()`: reads `CNTVCT_EL0`, stores as base ticks. Calls `time.Since(epoch).Milliseconds()` once for base milliseconds.
2. `Now()`: reads `CNTVCT_EL0`, subtracts base ticks, divides by ticks-per-ms, adds base milliseconds.
3. `ticksPerMs` computed once at init from `CNTFRQ_EL0` (24 MHz on Apple Silicon = 24,000 ticks/ms).

**Fallback path (!arm64):**
- Uses `time.Since(epoch).Milliseconds()` directly.

## Performance

Results on a MacBook Pro M4 Max (Go 1.26, darwin/arm64):

```
BenchmarkNow-16        685477630         1.761 ns/op        0 B/op    0 allocs/op
BenchmarkTimeSince-16  100000000        11.94 ns/op         0 B/op    0 allocs/op
```

**6.7x faster** than `time.Since().Milliseconds()` on macOS arm64. On other
architectures both benchmarks produce identical results (fallback path).

## Requirements

- Go 1.26+
- arm64 fast path requires native arm64 or standard VM (no trapped `CNTVCT_EL0`)

## License

Apache-2.0 -- see [LICENSE](LICENSE).
