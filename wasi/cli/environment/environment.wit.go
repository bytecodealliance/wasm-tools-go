package environment

import "github.com/ydnar/wasm-tools-go/cm"

// GetEnvironment represents the imported function "wasi:cli/environment.get-environment".
//
// Get the POSIX-style environment variables.
//
// Each environment variable is provided as a pair of string variable names
// and string value.
//
// Morally, these are a value import, but until value imports are available
// in the component model, this import function should return the same
// values each time it is called.
func GetEnvironment() cm.List[cm.Tuple[string, string]] {
	return get_environment()
}

//go:wasmimport wasi:cli/environment@0.2.0-rc-2023-12-05 get-environment
func get_environment() cm.List[cm.Tuple[string, string]]

// GetArguments represents the imported function "wasi:cli/environment.get-arguments".
//
// Get the POSIX-style arguments to the program.
func GetArguments() cm.List[string] {
	return get_arguments()
}

//go:wasmimport wasi:cli/environment@0.2.0-rc-2023-12-05 get-arguments
func get_arguments() cm.List[string]

// InitialCWD represents the imported function "wasi:cli/environment.initial-cwd".
//
// Return a path that programs should use as their initial current working
// directory, interpreting `.` as shorthand for this.
func InitialCWD() cm.Option[string] {
	return initial_cwd()
}

//go:wasmimport wasi:cli/environment@0.2.0-rc-2023-12-05 initial-cwd
func initial_cwd() cm.Option[string]
