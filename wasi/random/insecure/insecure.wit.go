// Package insecure represents the interface "wasi:random/insecure".
//
// The insecure interface for insecure pseudo-random numbers.
//
// It is intended to be portable at least between Unix-family platforms and
// Windows.
package insecure

import "github.com/ydnar/wasm-tools-go/cm"

// GetInsecureRandomBytes represents the imported function "wasi:random/insecure.get-insecure-random-bytes".
//
// Return `len` insecure pseudo-random bytes.
//
// This function is not cryptographically secure. Do not use it for
// anything related to security.
//
// There are no requirements on the values of the returned bytes, however
// implementations are encouraged to return evenly distributed values with
// a long period.
func GetInsecureRandomBytes(len uint64) (result cm.List[uint8]) {
	get_insecure_random_bytes(len, &result)
	return
}

//go:wasmimport wasi:random/insecure@0.2.0 get-insecure-random-bytes
func get_insecure_random_bytes(len uint64, result *cm.List[uint8])

// GetInsecureRandomU64 represents the imported function "wasi:random/insecure.get-insecure-random-u64".
//
// Return an insecure pseudo-random `u64` value.
//
// This function returns the same type of pseudo-random data as
// `get-insecure-random-bytes`, represented as a `u64`.
func GetInsecureRandomU64() uint64 {
	return get_insecure_random_u64()
}

//go:wasmimport wasi:random/insecure@0.2.0 get-insecure-random-u64
func get_insecure_random_u64() uint64
