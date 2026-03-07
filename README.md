# mclock

Fast monotonic clock for Go with millisecond and microsecond resolution.

On darwin/arm64 (Apple Silicon), reads the `CNTVCT_EL0` counter register directly via a single `MRS` instruction (~1.8 ns), bypassing Go's runtime and libSystem. Falls back to `time.Since` on all other platforms.

## Usage

```go
import "github.com/bds421/rho-mclock"

epoch := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
clk := mclock.New(epoch)
ms := clk.Now()      // milliseconds since epoch
us := clk.NowMicro() // microseconds since epoch
```

## How It Works

**darwin/arm64 fast path:**
1. `New()`: reads `CNTVCT_EL0`, stores as base ticks. Samples `time.Since(epoch)` once for base milliseconds and microseconds.
2. `Now()`: reads `CNTVCT_EL0`, subtracts base ticks, converts to milliseconds using exact quotient+remainder arithmetic, adds base milliseconds.
3. `NowMicro()`: same as `Now()` but converts to microseconds.
4. The counter frequency is read once at init from `CNTFRQ_EL0` (24 MHz on Apple Silicon). Conversion uses `delta/freq*scale + delta%freq*scale/freq` to avoid the drift that integer-truncated divisors introduce on frequencies that are not exact multiples of 1000.

**Fallback path:**
- Used on all non-darwin/arm64 platforms, and on darwin/arm64 if `CNTFRQ_EL0` returns 0.
- `New()` re-anchors the epoch with a monotonic reading from `time.Now()`.
- `Now()` uses `time.Since(epoch).Milliseconds()` directly.
- `NowMicro()` uses `time.Since(epoch).Microseconds()` directly.

## Performance

Results are environment-specific. MacBook Pro M4 Max, Go 1.26, darwin/arm64:

| Benchmark | ns/op | vs stdlib |
|---|---|---|
| `mclock.Now()` | 1.76 | **6.7x faster** |
| `time.Since().Milliseconds()` | 11.94 | baseline |

Zero allocations for both. On non-darwin/arm64 platforms both paths are equivalent.

## Requirements

- Go 1.26+
- Fast path requires darwin/arm64 (Apple Silicon). Linux arm64 and all other platforms use the fallback path.
- The fast path reads `CNTVCT_EL0` directly; VMs or environments that trap this register will get the fallback path automatically.

## License

Apache-2.0 -- see [LICENSE](LICENSE).
