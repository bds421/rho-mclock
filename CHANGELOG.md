# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.1] - 2026-03-14

### Changed
- Bump Go version to 1.26.1
- Add `#nosec G115` annotations on tick conversion functions (elapsed time overflow is physically impossible)

### Fixed
- Widen `TestNowAgreesWithTimeSince` tolerance from ±1ms to ±2ms

## [0.2.0] - 2026-03-07

### Fixed
- Fix conversion drift in tick-to-duration calculation
- Remove `init()` panics; errors are now returned from constructors

### Changed
- Tighten platform scope: fast path restricted to `darwin/arm64`

## [0.1.0] - 2026-03-07

### Added
- `Clock` interface with `Now()` (milliseconds) and `NowMicro()` (microseconds)
- Fast path for `darwin/arm64`: reads `CNTVCT_EL0` hardware counter directly (~1.76 ns/op)
- Fallback path for all other platforms using `time.Since` (~11.94 ns/op)
- `New()` constructor returning a `Clock` instance
- Comprehensive test suite with race-condition detection
- Benchmarks with `b.Loop()` (Go 1.24+)
- `Makefile` with `test`, `test-race`, `bench`, `vet`, and `clean` targets
- Apache 2.0 license
