# mclock

Fast monotonic clock for Go with millisecond and microsecond resolution.

On arm64, reads the `CNTVCT_EL0` counter register directly via a single `MRS` instruction (~1.8 ns), bypassing Go's runtime and libSystem. Falls back to `time.Since` on other architectures.

## Usage

```go
import "github.com/bds421/rho-mclock"

epoch := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
clk := mclock.New(epoch)
ms := clk.Now()      // milliseconds since epoch
us := clk.NowMicro() // microseconds since epoch
```

## How It Works

**arm64 fast path:**
1. `New()`: reads `CNTVCT_EL0`, stores as base ticks. Samples `time.Since(epoch)` once for base milliseconds and microseconds.
2. `Now()`: reads `CNTVCT_EL0`, subtracts base ticks, divides by ticks-per-ms, adds base milliseconds.
3. `NowMicro()`: same as `Now()` but divides by ticks-per-us and adds base microseconds.
4. `ticksPerMs` and `ticksPerUs` computed once at init from `CNTFRQ_EL0` (24 MHz on Apple Silicon = 24,000 ticks/ms, 24 ticks/us).

**Fallback path (!arm64):**
- `Now()` uses `time.Since(epoch).Milliseconds()` directly.
- `NowMicro()` uses `time.Since(epoch).Microseconds()` directly.

## Performance

MacBook Pro M4 Max, Go 1.26, darwin/arm64:

| Benchmark | ns/op | vs stdlib |
|---|---|---|
| `mclock.Now()` | 1.76 | **6.7x faster** |
| `time.Since().Milliseconds()` | 11.94 | baseline |

Zero allocations for both. On non-arm64 architectures both paths are equivalent.

## Requirements

- Go 1.26+
- arm64 fast path requires native arm64 or standard VM (no trapped `CNTVCT_EL0`)

## License

Apache-2.0 -- see [LICENSE](LICENSE).
