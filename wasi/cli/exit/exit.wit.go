package exit

import "github.com/ydnar/wasm-tools-go/cm"

// Exit represents the imported function "wasi:cli/exit.exit".
//
// Exit the current instance and any linked instances.
func Exit(status cm.UntypedResult) {
	exit(status)
}

//go:wasmimport wasi:cli/exit@0.2.0 exit
func exit(status cm.UntypedResult)
