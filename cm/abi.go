package cm

import (
	"unsafe"
)

// LowerResult lowers an untyped result into Core WebAssembly I32.
func LowerResult[T ~bool](v T) uint32 {
	return uint32(*(*uint8)(unsafe.Pointer(&v)))
}

// LowerHandle lowers a handle ([cm.Resource], [cm.Rep]) into a Core WebAssembly I32.
func LowerHandle[T any](v T) uint32 {
	return *(*uint32)(unsafe.Pointer(&v))
}

// LowerEnum lowers an enum into a Core WebAssembly I32.
func LowerEnum[E ~uint8 | ~uint16 | ~uint32](e E) uint32 {
	return uint32(e)
}

// LowerString lowers a [string] into a pair of Core WebAssembly types.
func LowerString[S ~string](s S) (*byte, uint) {
	return unsafe.StringData(string(s)), uint(len(s))
}

// LowerList lowers a [List] into a pair of Core WebAssembly types.
func LowerList[T any](list List[T]) (*T, uint) {
	return list.Data(), list.Len()
}

// LiftList lifts Core WebAssembly types into a [List].
func LiftList[L List[T], T any, Data unsafe.Pointer | uintptr | *T, Len uint | uintptr | uint32 | uint64](data Data, len Len) L {
	return L{
		data: (*T)(unsafe.Pointer(data)),
		len:  uint(len),
	}
}

func LowerBool[B ~bool](v B) uint32 { return uint32(*(*uint8)(unsafe.Pointer(&v))) }

func BoolToU32[B ~bool](v B) uint32 { return uint32(*(*uint8)(unsafe.Pointer(&v))) }
func BoolToU64[B ~bool](v B) uint64 { return uint64(*(*uint8)(unsafe.Pointer(&v))) }
func S32ToF32(v int32) float32      { return *(*float32)(unsafe.Pointer(&v)) }
func S64ToF64(v int64) float64      { return *(*float64)(unsafe.Pointer(&v)) }
func F64ToU64(v float64) uint64     { return *(*uint64)(unsafe.Pointer(&v)) }
func U32ToF32(v uint32) float32     { return *(*float32)(unsafe.Pointer(&v)) }
func U64ToF64(v uint64) float64     { return *(*float64)(unsafe.Pointer(&v)) }
func F32ToF64(v float32) float64    { return float64(v) }

func PointerToU32[T any](v *T) uint32 { return uint32(uintptr(unsafe.Pointer(v))) }
func PointerToU64[T any](v *T) uint64 { return uint64(uintptr(unsafe.Pointer(v))) }
