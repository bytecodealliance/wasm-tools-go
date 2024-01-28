package stderr

import "github.com/ydnar/wasm-tools-go/wasi/io/streams"

type OutputStream = streams.OutputStream

// GetStderr represents the imported function "wasi:cli/stderr#get-stderr".
func GetStderr() OutputStream {
	return get_stderr()
}

//go:wasmimport wasi:cli/stderr@0.2.0 get-stderr
func get_stderr() OutputStream
