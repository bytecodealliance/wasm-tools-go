// Package random represents the interface "wasi:random".
//
// WASI Random is a random data API.
//
// It is intended to be portable at least between Unix-family platforms and
// Windows.

package random

import "github.com/ydnar/wasm-tools-go/cm"

// GetRandomBytes represents the imported function "wasi:random/random.get-random-bytes".
//
// Return `len` cryptographically-secure random or pseudo-random bytes.
//
// This function must produce data at least as cryptographically secure and
// fast as an adequately seeded cryptographically-secure pseudo-random
// number generator (CSPRNG). It must not block, from the perspective of
// the calling program, under any circumstances, including on the first
// request and on requests for numbers of bytes. The returned data must
// always be unpredictable.
//
// This function must always return fresh data. Deterministic environments
// must omit this function, rather than implementing it with deterministic
// data.
func GetRandomBytes(len uint64) (result cm.List[uint8]) {
	get_random_bytes(len, &result)
	return
}

//go:wasmimport wasi:random/random@0.2.0 get-random-bytes
func get_random_bytes(len uint64, result *cm.List[uint8])

// GetRandomU64 represents the imported function "wasi:random/random.get-random-u64".
//
// Return a cryptographically-secure random or pseudo-random `u64` value.
//
// This function returns the same type of data as `get-random-bytes`,
// represented as a `u64`.
func GetRandomU64() uint64 {
	return get_random_u64()
}

//go:wasmimport wasi:random/random@0.2.0 get-random-u64
func get_random_u64() uint64
