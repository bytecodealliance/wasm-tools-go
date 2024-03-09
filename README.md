# wasm-tools-go

[![build status](https://img.shields.io/github/actions/workflow/status/ydnar/wasm-tools-go/test.yaml?branch=main)](https://github.com/ydnar/wasm-tools-go/actions)
[![pkg.go.dev](https://img.shields.io/badge/docs-pkg.go.dev-blue.svg)](https://pkg.go.dev/github.com/ydnar/wasm-tools-go)

## [WebAssembly](https://webassembly.org) + [WASI](https://wasi.dev) tools for [Go](https://go.dev)

## About

This repository contains code to generate Go bindings for [Component Model](https://component-model.bytecodealliance.org/) interfaces defined in [WIT](https://component-model.bytecodealliance.org/design/wit.html) (WebAssembly Interface Type) files. The goal of this project is to accelerate adoption of the Component Model and development of [WASI Preview 2](https://bytecodealliance.org/articles/WASI-0.2) in Go.

### WASI

The [wasi](./wasi) directory contains generated bindings for the [`wasi:cli/command@0.2.0`](https://github.com/WebAssembly/wasi-cli) world from [WASI Preview 2](https://github.com/WebAssembly/WASI/blob/main/preview2/README.md). The generated bindings currently target [TinyGo](https://tinygo.org).

### Component Model

Package [cm](./cm) contains low-level types used by the `wasi` packages, such as `option<t>`, `result<ok, err>`, `variant`, `list`, and `resource`. These are intended for use by generated [Component Model](https://github.com/WebAssembly/component-model/blob/main/design/mvp/Explainer.md#type-definitions) bindings, where the caller converts to a Go equivalent. It attempts to map WIT semantics to their equivalent in Go where possible.

#### Note on Memory Safety

Package `cm` and generated bindings from `wit-bindgen-go` may have compatibility issues with the Go garbage collector, as they directly represent `variant` and `result` types as tagged unions where a pointer shape may be occupied by a non-pointer value. The GC may detect and throw an error if it detects a non-pointer value in an area it expects to see a pointer. This is an area of active development.

## `wit-bindgen-go`

### WIT → Go

The `wit-bindgen-go` tool can generate Go bindings for WIT interfaces and worlds. If [`wasm-tools`](https://crates.io/crates/wasm-tools) is installed and in `$PATH`, then `wit-bindgen-go` can load WIT directly.

```sh
wit-bindgen-go generate ../wasi-cli/wit
```

Otherwise, pass the JSON representation of a fully-resolved WIT package:

```sh
wit-bindgen-go generate wasi-cli.wit.json
```

Or pipe via `stdin`:

```sh
wasm-tools component wit -j ../wasi-cli/wit | wit-bindgen-go generate
```

### JSON → WIT

For debugging purposes, `wit-bindgen-go` can also convert a JSON representation back into WIT. This is useful for validating that the intermediate representation faithfully represents the original WIT source.

```sh
wit-bindgen-go wit example.wit.json
```

### WIT → JSON

The [wit](./wit) package can decode a JSON representation of a fully-resolved WIT file. Serializing WIT into JSON requires [wasm-tools](https://crates.io/crates/wasm-tools) v1.0.42 or higher. To convert a WIT file into JSON, run `wasm-tools` with the `-j` argument:

```sh
wasm-tools component wit -j example.wit
```

This will emit JSON on `stdout`, which can be piped to a file or another program.

```sh
wasm-tools component wit -j example.wit > example.wit.json
```

## Notes

### Export Bindings for Resource Types

For each exported resource type, in addition to constructors, methods, and static functions, a component must import and export a number of functions to manage the resource lifecycle.

### Example

For a hypothetical WIT resource `water` in package/interface `example:foo/bar`, a component _imports_ the following administrative functions:

- **resource-new**: initialize a resource handle for a given `rep` (concrete representation). Called by the component before a resource value is returned by an exported function.
	- Import: `[export]example:foo/bar [resource-new]water`
	- Params: `i32` rep (in C bindings, this is a pointer to the representation)
	- Results: `i32` handle
- **resource-rep**: return the underlying representation of a handle. Can only be called by the component instance that created the original resource.
	- Import: `[export]example:foo/bar [resource-rep]water`
	- Params: `i32` handle
	- Results: `i32` rep (in C bindings, this is a pointer to the representation)
- **resource-drop**: drops an owned handle to a resource. If the resource has zero loans, the destructor (below) will be called.
	- Import: `example:foo/bar [resource-drop]water`
	- Params: `i32` handle
	- Results: none

In addition, the resource destructor is _exported_:

- **dtor**: destructor, implemented by user code
	- Export: `example:foo/bar#[dtor]water`
	- Params: `i32` rep (in C bindings, this is a pointer)
	- Results: none
	- Notes: the

Similarly, a function is exported for each resource method (e.g. `drink` and `spill`)

- Export: `example:foo/bar#[method]water.drink` (implemented by user code)
- Export: `example:foo/bar#[method]water.spill` (implemented by user code)

## License

This project is licensed under the Apache 2.0 license with the LLVM exception. See [LICENSE](LICENSE) for more details.

### Contributing

Developing and testing `wit-bindgen-go` requires an up-to-date installation of [Go](https://go.dev), [TinyGo](https://tinygo.org), and [`wasm-tools`](https://crates.io/crates/wasm-tools).

Unless you explicitly state otherwise, any contribution intentionally submitted for inclusion in this project by you, as defined in the Apache-2.0 license, shall be licensed as above, without any additional terms or conditions.
