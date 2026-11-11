# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Fixed

- `Must` panics with the original error value so identity survives `recover`.
- Wrapping preserves the original non-empty code across nested / non-errstack wrappers (`CodeOf`).

### Added

- Structured `Error` type with `Code`, `Message`, and `StackTrace`.
- Independent `Frame` / `StackTrace` types with exact human-readable `String`.
- Stack capture via `runtime.Callers` and `runtime.CallersFrames`.
- `FormatStackTrace`, `FormatText`, and `fmt.Formatter` (`%+v`) support.
- Structured JSON serialization of `stack_trace` as an array of frames.
- `NewCode` / `NewCodef` constructors.
- `Trace` helper returning a cloned structured stack.
- Wrap/WithStack preserve an existing errstack stack instead of replacing it.

### Changed

- Replaced `github.com/pkg/errors` with a stdlib-only implementation.
- `Stack` / `ErrorWithStack` now format the structured stack model.

### Notes

- HTTP status codes remain intentionally out of scope for this library.

## [0.1.0] - 2026-11-11

### Added

- Initial extraction of stacked error helpers from `Onyx.Common/errstack`.
- Create/wrap/cause/stack helpers built on `github.com/pkg/errors`.

[Unreleased]: https://github.com/MahdiPezeshkian/onyx-errors/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/MahdiPezeshkian/onyx-errors/releases/tag/v0.1.0
