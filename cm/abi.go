package cm

import (
	"unsafe"
)

// Reinterpret reinterprets the bits of type From into type To.
// Will panic if the size of From is smaller than To.
func Reinterpret[To, From any](from From) (to To) {
	if unsafe.Sizeof(to) > unsafe.Sizeof(from) {
		panic("reinterpret: size of to > from")
	}
	return *(*To)(unsafe.Pointer(&from))
}

// LowerResult lowers an untyped result into Core WebAssembly I32.
func LowerResult[T ~bool](v T) uint32 {
	return uint32(*(*uint8)(unsafe.Pointer(&v)))
}

// LowerHandle lowers a handle ([cm.Resource], [cm.Rep]) into a Core WebAssembly I32.
func LowerHandle[T any](v T) uint32 {
	return *(*uint32)(unsafe.Pointer(&v))
}

// LiftHandle lifts Core WebAssembly I32 into a handle ([cm.Resource], [cm.Rep]).
func LiftHandle[H any](v uint32) H {
	return *(*H)(unsafe.Pointer(&v))
}

// LowerEnum lowers an enum into a Core WebAssembly I32.
func LowerEnum[E ~uint8 | ~uint16 | ~uint32](e E) uint32 {
	return uint32(e)
}

// LowerString lowers a [string] into a pair of Core WebAssembly types.
func LowerString[S ~string](s S) (*byte, uint) {
	return unsafe.StringData(string(s)), uint(len(s))
}

// LiftString lifts Core WebAssembly types into a [string].
func LiftString[Data unsafe.Pointer | uintptr | *uint8, Len uint | uintptr | uint32 | uint64](data Data, len Len) string {
	return unsafe.String((*uint8)(unsafe.Pointer(data)), int(len))
}

// LowerList lowers a [List] into a pair of Core WebAssembly types.
func LowerList[L ~struct {
	data *T
	len  uint
}, T any](list L) (*T, uint) {
	l := (*List[T])(unsafe.Pointer(&list))
	return l.data, l.len
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
