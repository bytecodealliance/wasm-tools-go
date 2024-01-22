package stdout

import "github.com/ydnar/wasm-tools-go/wasi/io/streams"

type OutputStream = streams.OutputStream

// GetStdout represents the imported function "wasi:cli/stdout#get-stdout".
func GetStdout() OutputStream {
	return get_stdout()
}

//go:wasmimport wasi:cli/stdout@0.2.0-rc-2023-11-10 get-stdout
func get_stdout() OutputStream
