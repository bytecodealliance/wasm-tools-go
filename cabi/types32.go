//go:build wasm32 || tinygo.wasm

package cabi

import "unsafe"

type alloc32 = unsafe.Pointer
type alloc64 = [2]unsafe.Pointer
type alloc128 = [4]unsafe.Pointer

func init() {
	var a32 alloc32
	if unsafe.Sizeof(a32) != 4 {
		panic("sizeof alloc32 != 4")
	}
	var a64 alloc64
	if unsafe.Sizeof(a64) != 8 {
		panic("sizeof alloc64 != 8")
	}
	var a128 alloc128
	if unsafe.Sizeof(a128) != 16 {
		panic("sizeof alloc128 != 16")
	}
}
