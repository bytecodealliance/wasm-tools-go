package tinyunsafe

import "unsafe"

type anystruct any

// TODO: remove this when TinyGo supports unsafe.Offsetof
func OffsetOf[Struct anystruct, Field any](s *Struct, f *Field) uintptr {
	return uintptr(unsafe.Pointer(f)) - uintptr(unsafe.Pointer(s))
}
