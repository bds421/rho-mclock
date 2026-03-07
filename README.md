# mclock

Fast monotonic millisecond clock for Go.

On arm64, reads the `CNTVCT_EL0` counter register directly via a single `MRS` instruction (~5-8 ns), bypassing Go's runtime and libSystem. Falls back to `time.Since` on other architectures.

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

| Platform | `time.Since` | `mclock.Now` |
|---|---|---|
| macOS arm64 (M4 Max) | ~280 ns | ~5-8 ns |
| Linux arm64 | ~15-20 ns | ~5-8 ns |
| Linux/macOS amd64 | ~15-50 ns | same (fallback) |

## Requirements

- Go 1.26+
- arm64 fast path requires native arm64 or standard VM (no trapped `CNTVCT_EL0`)

## License

See [LICENSE](../rho-snowflake-2026/LICENSE).
