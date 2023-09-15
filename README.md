# wasm-tools-go

## [WebAssembly](https://webassembly.org) + [WASI](https://wasi.dev) tools for [Go](https://go.dev).

## About

This repository contains code to adapt [WIT](https://github.com/WebAssembly/component-model/blob/main/design/mvp/WIT.md) (Wasm Interface Type) files into Go, with the goal of accelerating the Go implementation of [WASI Preview 2](https://bytecodealliance.org/articles/webassembly-the-updated-roadmap-for-developers).

## WIT → JSON

This package can decode a JSON representation of a fully-resolved WIT file. Serializing WIT into JSON requires a version of wasm-tools as of commit [4975360](https://github.com/bytecodealliance/wasm-tools/commit/49753602683a539b66d0a65ffa11acb402f148bb) ([#1203](https://github.com/bytecodealliance/wasm-tools/pull/1203)), which adds serialization support to the [wit-parser](https://docs.rs/wit-parser/latest/wit_parser/) crate. To convert a WIT file into JSON, run `wasm-tools` with the `-j` argument:

```sh
wasm-tools component wit -j example.wit
```

This will emit JSON on stdout, which can be piped to a file or another program.

```sh
wasm-tools component wit -j example.wit > example.wit.json
```

```sh
wasm-tools component wit -j example.wit > wit-bindgen-go
```

## License

This project is licensed under the Apache 2.0 license with the LLVM exception. See [LICENSE](LICENSE) for more details.

### Contribution

Unless you explicitly state otherwise, any contribution intentionally submitted for inclusion in this project by you, as defined in the Apache-2.0 license, shall be licensed as above, without any additional terms or conditions.
