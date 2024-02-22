# wasm-tools-go

[![build status](https://img.shields.io/github/actions/workflow/status/ydnar/wasm-tools-go/test.yaml?branch=main)](https://github.com/ydnar/wasm-tools-go/actions)
[![pkg.go.dev](https://img.shields.io/badge/docs-pkg.go.dev-blue.svg)](https://pkg.go.dev/github.com/ydnar/wasm-tools-go)

## [WebAssembly](https://webassembly.org) + [WASI](https://wasi.dev) tools for [Go](https://go.dev)

## About

This repository contains code to adapt [WIT](https://github.com/WebAssembly/component-model/blob/main/design/mvp/WIT.md) (Wasm Interface Type) files into Go, with the goal of accelerating the Go implementation of [WASI Preview 2](https://bytecodealliance.org/articles/webassembly-the-updated-roadmap-for-developers).

## WIT → Go

The `wit-bindgen-go` tool can generate Go bindings for WIT interfaces and worlds. To use, pass the JSON representation of a fully-resolved WIT package:

```sh
wasm-tools component wit -j ./wasi-cli/wit | wit-bindgen-go generate
```

Or:

```sh
wit-bindgen-go generate wasi-cli.wit.json
```

## JSON → WIT

For debugging purposes, the `wit-bindgen-go` tool can also convert a JSON representation back into WIT. This is useful for validating that the intermediate representation faithfully represents the original WIT source.

```sh
wit-bindgen-go wit example.wit.json
```

## WIT → JSON

The [wit](./wit) package can decode a JSON representation of a fully-resolved WIT file. Serializing WIT into JSON requires [wasm-tools](https://crates.io/crates/wasm-tools) v1.0.42 or higher. To convert a WIT file into JSON, run `wasm-tools` with the `-j` argument:

```sh
wasm-tools component wit -j example.wit
```

This will emit JSON on `stdout`, which can be piped to a file or another program.

```sh
wasm-tools component wit -j example.wit > example.wit.json
```

## Go Packages

### WASI

The [wasi](./wasi) directory contains a machine-generated bindings for the [`wasi:cli/command`](https://github.com/WebAssembly/wasi-cli) world from [WASI Preview 2](https://github.com/WebAssembly/WASI/blob/main/preview2/README.md). The generated bindings currently target [TinyGo](https://tinygo.org).

### Component Model

Package [cm](./cm) contains low-level types used by the `wasi` packages, such as `option<t>`, `result<ok, err>`, `variant`, `list`, and `resource`. These are intended for use by generated [Component Model](https://github.com/WebAssembly/component-model/blob/main/design/mvp/Explainer.md#type-definitions) bindings, where the caller converts to a Go equivalent.

## License

This project is licensed under the Apache 2.0 license with the LLVM exception. See [LICENSE](LICENSE) for more details.

### Contribution

Unless you explicitly state otherwise, any contribution intentionally submitted for inclusion in this project by you, as defined in the Apache-2.0 license, shall be licensed as above, without any additional terms or conditions.
