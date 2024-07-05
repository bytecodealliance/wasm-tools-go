# Changelog

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Breaking Changes
- Package `cm`: Removed unused functions `Reinterpret2`, `LowerResult`, `LowerBool`, `BoolToU64`, `S64ToF64`.
- Package `cm`: Removed unused experimental `flags` implementation behind a build tag.

### Added
- Canonical ABI lifting code for `flags` and `variant` types.
- Lifting code for `result` and `variant` will now panic if caller passes an invalid case.
- Additional test coverage for `variant` and `flags` cases.

### Fixed
- Lifting code for `result` with no error type will now correctly set `IsErr`.

## [v0.1.1] — 2024-07-04

This release changes the memory layout of `variant` and `result` types to permit passing these types on the stack safely. This required breaking changes to package `cm`, detailed below, as well as slightly more verbose type signatures for WIT functions that return a typed `result`.

### Breaking Changes
- Type `cm.Result` is now `cm.BoolResult`.
- Types `cm.OKResult` and `cm.ErrResult` have been removed, replaced with a more generalized `cm.Result[Shape, OK, Err]` type.

### Added
- WIT labels with uppercase acronyms or [initialisms](https://go.dev/wiki/CodeReviewComments#initialisms) are now preserved in Go form. For example, the WIT name `time-EOD` will result in the Go name `TimeEOD`.
- `OK` is now a predefined initialism. For example, the WIT name `state-ok` would previously translate into the Go name `StateOk` instead of the idiomatic `StateOK`.

### Fixed
- [#95](https://github.com/ydnar/wasm-tools-go/issues/95): `wit-bindgen-go` now correctly generates packed `data` shape types for `variant` and `result` types.
- Fixed swapped `Shape` and `Align` type parameters in the functions `cm.New` and `cm.Case` for manipulating `variant` types.
- Variant validation now correctly reports `variant` instead of `result` in panic messages.

## [v0.1.0] — 2024-07-04

Initial version, supporting [TinyGo](https://tinygo.org/) + [WASI](https://wasi.dev) 0.2 (WASI Preview 2).

### Known Issues
- [#95](https://github.com/ydnar/wasm-tools-go/issues/95): `variant` and `result` types without fully-packed `data` shape types will not correctly represent all associated types.
- [#111](https://github.com/ydnar/wasm-tools-go/issues/111): `flags` types with > 32 labels are not correctly supported. See [component-model#370](https://github.com/WebAssembly/component-model/issues/370) and [wasm-tools#1635](https://github.com/bytecodealliance/wasm-tools/pull/1635) for more information.
- [#118](https://github.com/ydnar/wasm-tools-go/issues/118): Canonial ABI [post-return](https://github.com/WebAssembly/component-model/blob/main/design/mvp/CanonicalABI.md#canon-lift) functions to clean up allocations are not currently generated.
- Because Go does not have a native tagged union type, pointers represented in `variant` and `result` types may not be visible to the garbage collector and may be freed while still in use.
- Support for mainline [Go](https://go.dev/).

[Unreleased]: <https://github.com/ydnar/wasm-tools-go/compare/v0.1.1...HEAD>
[v0.1.1]: <https://github.com/ydnar/wasm-tools-go/compare/v0.1.0...v0.1.1>
[v0.1.0]: <https://github.com/ydnar/wasm-tools-go/tree/v0.1.0>
