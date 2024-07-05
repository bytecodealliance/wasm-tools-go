package cm

import (
	"unsafe"
)

// Reinterpret reinterprets the bits of type From into type T.
// Will panic if the size of From is smaller than the size of To.
func Reinterpret[T, From any](from From) (to T) {
	if unsafe.Sizeof(to) > unsafe.Sizeof(from) {
		panic("reinterpret: size of to > from")
	}
	return *(*T)(unsafe.Pointer(&from))
}

// LowerString lowers a [string] into a pair of Core WebAssembly types.
func LowerString[S ~string](s S) (*byte, uint32) {
	return unsafe.StringData(string(s)), uint32(len(s))
}

// LiftString lifts Core WebAssembly types into a [string].
func LiftString[T ~string, Data unsafe.Pointer | uintptr | *uint8, Len uint | uintptr | uint32 | uint64](data Data, len Len) T {
	return T(unsafe.String((*uint8)(unsafe.Pointer(data)), int(len)))
}

// LowerList lowers a [List] into a pair of Core WebAssembly types.
func LowerList[L ~struct{ list[T] }, T any](list L) (*T, uint32) {
	l := (*List[T])(unsafe.Pointer(&list))
	return l.data, uint32(l.len)
}

// LiftList lifts Core WebAssembly types into a [List].
func LiftList[L List[T], T any, Data unsafe.Pointer | uintptr | *T, Len uint | uintptr | uint32 | uint64](data Data, len Len) L {
	return L(NewList((*T)(unsafe.Pointer(data)), uint(len)))
}

func BoolToU32[B ~bool](v B) uint32 { return uint32(*(*uint8)(unsafe.Pointer(&v))) }
func BoolToU64[B ~bool](v B) uint64 { return uint64(*(*uint8)(unsafe.Pointer(&v))) }
func S64ToF64(v int64) float64      { return *(*float64)(unsafe.Pointer(&v)) }
func F64ToU64(v float64) uint64     { return *(*uint64)(unsafe.Pointer(&v)) }
func U32ToBool(v uint32) bool       { tmp := uint8(v); return *(*bool)(unsafe.Pointer(&tmp)) }
func U32ToF32(v uint32) float32     { return *(*float32)(unsafe.Pointer(&v)) }
func U64ToBool(v uint64) bool       { tmp := uint8(v); return *(*bool)(unsafe.Pointer(&tmp)) }
func U64ToF64(v uint64) float64     { return *(*float64)(unsafe.Pointer(&v)) }
func F32ToF64(v float32) float64    { return float64(v) }

func PointerToU32[T any](v *T) uint32 { return uint32(uintptr(unsafe.Pointer(v))) }
func PointerToU64[T any](v *T) uint64 { return uint64(uintptr(unsafe.Pointer(v))) }
func U32ToPointer[T any](v uint32) *T { return (*T)(unsafePointer(uintptr(v))) }
func U64ToPointer[T any](v uint64) *T { return (*T)(unsafePointer(uintptr(v))) }

// Appease vet, see https://github.com/golang/go/issues/58625
func unsafePointer(p uintptr) unsafe.Pointer {
	return *(*unsafe.Pointer)(unsafe.Pointer(&p))
}
