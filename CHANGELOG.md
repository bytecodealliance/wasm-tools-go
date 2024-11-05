# Changelog

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v0.3.2] — 2024-11-05

This module has been renamed to [go.bytecodealliance.org](https://go.bytecodealliance.org). Going forward, please use `go.bytecodealliance.org` instead of `github.com/bytecodealliance/wasm-tools-go`.

## [v0.3.1] — 2024-11-01

### Added

- `wit-bindgen-go` now accepts arguments that control the level of logging output on stderr. Verbose mode (`-v` or `--verbose`) will print informational log lines, warnings, and errors. Debug mode (`-vv` or `--debug`) will emit finer-grained debugging information. By default, `wit-bindgen-go` will print warnings or errors.
- New method `(cm.Result).Result() (ok OK, err Err, isErr bool)` added to streamline common patterns using `result` types. It has a value receiver, allowing it to be chained off a function call that returns a `Result`, eliminating the need to declare a temporary variable. For example: `stream, err, isErr := stream.BlockingRead(n)`

### Fixed

- [#215](https://github.com/bytecodealliance/wasm-tools-go/issues/215): generate variant accessor methods with the correct scope for Go names.

### New Contributors

- [@vados-cosmonic](https://github.com/vados-cosmonic): [#212](https://github.com/bytecodealliance/wasm-tools-go/issues/212) and [#214](https://github.com/bytecodealliance/wasm-tools-go/pull/214) — verbose logging controls in `wit-bindgen-go`.

## [v0.3.0] — 2024-10-11

### Added

- Go [type aliases](https://go.dev/ref/spec#Alias_declarations) are now generated for each WIT type alias (`type foo = bar`). Deep chains of type aliases (`type b = a; type c = b;`) are fully supported. Generated documentation now reflects whether a type is an alias. Fixes [#204](https://github.com/bytecodealliance/wasm-tools-go/issues/204).
- `go:wasmimport` and `go:wasmexport` functions are now generated in a separate `.wasm.go` file. This helps enable testing or use of generated packages outside of WebAssembly.
- `wit-bindgen-go generate` now generates a WIT file for each WIT world in its corresponding Go package directory. For example the `wasi:http/proxy` world would generate `wasi/http/proxy/proxy.wit`.
- `wit-bindgen-go wit` now accepts a `--world` argument in the form of `imports`, `wasi:clocks/imports`, or `wasi:clocks/imports@0.2.0`. This filters the serialized WIT to a specific world and interfaces it references. This can be used to generate focused WIT for a specific world with a minimal set of dependencies.

### Changed

- Method `wit.(*Package).WIT()` now interprets the non-empty string `name` argument as signal to render in single-file, multi-package braced form.
- `wit.(*Resolve).WIT()` and `wit.(*Package).WIT()` now accept a `*wit.World` as context to filter serialized WIT to a specific world.
- Packages are now sorted topologically by dependency in generated WIT files. For example, `wasi:cli` would be followed by its dependencies like `wasi:io`, `wasi:filesystem`, or `wasi:random`.

## [v0.2.4] — 2024-10-06

### Added

- Generated variant shape types (`...Shape`) now include [`structs.HostLayout`](https://github.com/golang/go/issues/66984).

## [v0.2.3] — 2024-10-05

### Added

- `wit-bindgen-go generate` now accepts a remote registry reference to pull down a WIT package and generate the Go code. Example: `wit-bindgen-go generate ghcr.io/webassembly/wasi/http:0.2.0`.

### Changed

- `cm.List` now stores list length as a `uintptr`, permitted by the [Go wasm types proposal](https://github.com/golang/go/issues/66984). It was previously a `uint`, which was removed from the list of permitted types. There should be no change in the memory layout in TinyGo `GOARCH=wasm` or Go `GOARCH=wasm32` using 32-bit pointers.
- The helper functions `cm.NewList`, `cm.LiftList`, and `cm.LiftString` now accept any integer type for `len`, defined as `cm.AnyInteger`.
- `cm.Option[T].Value()` method now value receiver (not pointer receiver), so it can be chained.

## [v0.2.2] — 2024-10-03

### Added

- All `struct` types in package `cm` now include `structs.HostLayout` on Go 1.23 or later.
- Added type constraints `AnyList`, `AnyResult`, and `AnyVariant` in package `cm` to constrain generic functions accepting any of those types.
- Variant types now implement `fmt.Stringer`, with a `String` method. Breaking: variant cases named `string` map to `String_` in Go.
- `cm.Option[T]` types now have a `Value()` convenience method that returns the zero value for `T` if the option represents the none case. For example, this simplifies getting an empty string or slice from `option<string>` or `option<list<T>>`, respectively.
- Added a release workflow to publish tagged releases to GitHub.

### Fixed

- `wit-bindgen-go --version` now displays the version without empty `()`.

## [v0.2.1] — 2024-09-26

### Added

- Generated structs and structs in package `cm` now include a [`HostLayout` field](https://github.com/golang/go/issues/66408) in order to conform with the [relaxed types proposal](https://github.com/golang/go/issues/66984) for `GOARCH=wasm32`. The `cm.HostLayout` type is an alias for `structs.HostLayout` on Go 1.23 or later, and a polyfill for Go 1.22 or earlier.
- [#163](https://github.com/bytecodealliance/wasm-tools-go/issues/163): added `cm.F32ToU64()` and `cm.U64ToF32()` for flattening `f32` and `u64` types in the Canonical ABI.
- Test data from [bytecodealliance/wit-bindgen/tests/codegen](https://github.com/bytecodealliance/wit-bindgen/tree/main/tests/codegen).

### Fixed

- [#159](https://github.com/bytecodealliance/wasm-tools-go/pull/159): correctly escape all WIT keywords, including when used in package names.
- [#160](https://github.com/bytecodealliance/wasm-tools-go/issues/160): fixed the use of Go reserved keywords in function returns as result types.
- [#161](https://github.com/bytecodealliance/wasm-tools-go/issues/161): correctly handle `constructor` as a WIT keyword in wit.
- [#165](https://github.com/bytecodealliance/wasm-tools-go/issues/165): fixed use of imported types in exported functions.
- [#167](https://github.com/bytecodealliance/wasm-tools-go/issues/167): fixed a logic flaw in `TestHasBorrow`.
- [#170](https://github.com/bytecodealliance/wasm-tools-go/issues/170): resolve implied names for interface imports and exports in a world.
- [#175](https://github.com/bytecodealliance/wasm-tools-go/issues/175): generated correct symbol names for imported and exported functions in worlds (`$root`) or interfaces declared inline in worlds.

## [v0.2.0] — 2024-09-05

**This project has moved!** `wasm-tools-go` is now an official [Bytecode Alliance](https://github.com/bytecodealliance) project.

Going forward, please update your Go imports from `github.com/ydnar/wasm-tools-go` to `github.com/bytecodealliance/wasm-tools-go`. Thanks to [@ricochet](https://github.com/ricochet), [@mossaka](https://github.com/mossaka), [@lxfontes](https://github.com/lxfontes), and others for their help making this possible.

### Changed

- Added support for `@deprecated` directive implemented in [`wasm-tools#1687`](https://github.com/bytecodealliance/wasm-tools/pull/1687).
- Removed support for `@since` feature gating implemented in [`wasm-tools#1741`](https://github.com/bytecodealliance/wasm-tools/pull/1741).

### Fixed

- [#151](https://github.com/bytecodealliance/wasm-tools-go/issues/151): backport support for JSON generated by `wasm-tools` prior to v1.209.0, which added `@since` and `@unstable` feature gates.

## [v0.1.5] — 2024-08-23

### Added

- `wit-bindgen-go --version` now reports the module version and git revision of the `wit-bindgen-go` command.

### Changed

- Omit the default `//go:build !wasip1` build tags from generated Go files. This enables `wit-bindgen-go` to target `GOOS=wasip1` (fixes [#147](https://github.com/bytecodealliance/wasm-tools-go/issues/147)).
- Package `wit` now serializes multi-package WIT files with an un-nested “root” package. See [WebAssembly/component-model#380](https://github.com/WebAssembly/component-model/pull/380) and [bytecodealliance/wasm-tools#1700](https://github.com/bytecodealliance/wasm-tools/pull/1700).

## [v0.1.4] — 2024-07-16

### Added

- `wit-bindgen-go generate` now accepts a `--cm` option to specify the Go import path to package `cm`. Used for custom or internal implementations of package `cm`. Defaults to `github.com/bytecodealliance/wasm-tools-go/cm`.
- `Tuple9`...`Tuple16` types in package `cm` to align with [component-model#373](https://github.com/WebAssembly/component-model/issues/373). Tuples with 9 to 16 types will no longer generate inline `struct` types.
- Documentation for Canonical ABI lift and lower helper functions in package `cm`.

### Changed

- Removed outdated documentation in [design](./design/README.md).

## [v0.1.3] — 2024-07-08

### Added

- [#128](https://github.com/bytecodealliance/wasm-tools-go/pull/128): implemented `String` method for `enum` types ([@rajatjindal](https://github.com/rajatjindal)).

### Fixed

- [#130](https://github.com/bytecodealliance/wasm-tools-go/issues/130): anonymous `tuple` types now correctly have exported Go struct fields.
- [#129](https://github.com/bytecodealliance/wasm-tools-go/issues/129): correctly handle zero-length `tuple` and `record` types, represented as `struct{}`.

## [v0.1.2] — 2024-07-05

### Added

- Canonical ABI lifting code for `flags` and `variant` types.
- Lifting code for `result` and `variant` will now panic if caller passes an invalid case.
- Additional test coverage for `variant` and `flags` cases.

### Removed

- Package `cm`: Removed unused functions `Reinterpret2`, `LowerResult`, `LowerBool`, `BoolToU64`, `S64ToF64`.
- Package `cm`: Removed unused experimental `flags` implementation behind a build tag.

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

- [#95](https://github.com/bytecodealliance/wasm-tools-go/issues/95): `wit-bindgen-go` now correctly generates packed `data` shape types for `variant` and `result` types.
- Fixed swapped `Shape` and `Align` type parameters in the functions `cm.New` and `cm.Case` for manipulating `variant` types.
- Variant validation now correctly reports `variant` instead of `result` in panic messages.

## [v0.1.0] — 2024-07-04

Initial version, supporting [TinyGo](https://tinygo.org/) + [WASI](https://wasi.dev) 0.2 (WASI Preview 2).

### Known Issues

- [#95](https://github.com/bytecodealliance/wasm-tools-go/issues/95): `variant` and `result` types without fully-packed `data` shape types will not correctly represent all associated types.
- [#111](https://github.com/bytecodealliance/wasm-tools-go/issues/111): `flags` types with > 32 labels are not correctly supported. See [component-model#370](https://github.com/WebAssembly/component-model/issues/370) and [wasm-tools#1635](https://github.com/bytecodealliance/wasm-tools/pull/1635) for more information.
- [#118](https://github.com/bytecodealliance/wasm-tools-go/issues/118): Canonial ABI [post-return](https://github.com/WebAssembly/component-model/blob/main/design/mvp/CanonicalABI.md#canon-lift) functions to clean up allocations are not currently generated.
- Because Go does not have a native tagged union type, pointers represented in `variant` and `result` types may not be visible to the garbage collector and may be freed while still in use.
- Support for mainline [Go](https://go.dev/).

[Unreleased]: <https://github.com/bytecodealliance/wasm-tools-go/compare/v0.3.2..HEAD>
[v0.3.2]: <https://github.com/bytecodealliance/wasm-tools-go/compare/v0.3.1..v0.3.2>
[v0.3.1]: <https://github.com/bytecodealliance/wasm-tools-go/compare/v0.3.0..v0.3.1>
[v0.3.0]: <https://github.com/bytecodealliance/wasm-tools-go/compare/v0.2.4..v0.3.0>
[v0.2.4]: <https://github.com/bytecodealliance/wasm-tools-go/compare/v0.2.3..v0.2.4>
[v0.2.3]: <https://github.com/bytecodealliance/wasm-tools-go/compare/v0.2.2..v0.2.3>
[v0.2.2]: <https://github.com/bytecodealliance/wasm-tools-go/compare/v0.2.1..v0.2.2>
[v0.2.1]: <https://github.com/bytecodealliance/wasm-tools-go/compare/v0.2.0...v0.2.1>
[v0.2.0]: <https://github.com/bytecodealliance/wasm-tools-go/compare/v0.1.5...v0.2.0>
[v0.1.5]: <https://github.com/bytecodealliance/wasm-tools-go/compare/v0.1.4...v0.1.5>
[v0.1.4]: <https://github.com/bytecodealliance/wasm-tools-go/compare/v0.1.3...v0.1.4>
[v0.1.3]: <https://github.com/bytecodealliance/wasm-tools-go/compare/v0.1.2...v0.1.3>
[v0.1.2]: <https://github.com/bytecodealliance/wasm-tools-go/compare/v0.1.1...v0.1.2>
[v0.1.1]: <https://github.com/bytecodealliance/wasm-tools-go/compare/v0.1.0...v0.1.1>
[v0.1.0]: <https://github.com/bytecodealliance/wasm-tools-go/tree/v0.1.0>
