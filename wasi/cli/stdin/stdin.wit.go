package stdin

import "github.com/ydnar/wasm-tools-go/wasi/io/streams"

type InputStream = streams.InputStream

// GetStdin represents the imported function "wasi:cli/stdin#get-stdin".
func GetStdin() InputStream {
	return get_stdin()
}

//go:wasmimport wasi:cli/stdin@0.2.0 get-stdin
func get_stdin() InputStream
