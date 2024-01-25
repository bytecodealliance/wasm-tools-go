package preopens

import (
	"github.com/ydnar/wasm-tools-go/cm"
	"github.com/ydnar/wasm-tools-go/wasi/filesystem/types"
)

type Descriptor = types.Descriptor

// GetDirectories represents the imported function "wasi:filesystem/preopens.get-directories".
//
// Return the set of preopened directories, and their path.
func GetDirectories() (result cm.List[cm.Tuple[Descriptor, string]]) {
	get_directories(&result)
	return
}

//go:wasmimport wasi:filesystem/preopens@0.2.0-rc-2023-11-10 get-directories
func get_directories(result *cm.List[cm.Tuple[Descriptor, string]])
