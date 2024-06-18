package cm

import (
	"unsafe"
)

type CoreTypes interface {
	uint32 | uint64 | float32 | float64
}

type CorePointers[T any] interface {
	*T | unsafe.Pointer | uintptr
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

// Lower functions
func BoolToS32(v bool) int32  { return int32(*(*int8)(unsafe.Pointer(&v))) }
func BoolToU32(v bool) uint32 { return uint32(*(*uint8)(unsafe.Pointer(&v))) }
func BoolToS64(v bool) int64  { return int64(*(*int8)(unsafe.Pointer(&v))) }
func BoolToU64(v bool) uint64 { return uint64(*(*uint8)(unsafe.Pointer(&v))) }

// func BoolToF32(v bool) float32 { return U32ToF32(BoolToU32(v)) }
// func BoolToF64(v bool) float64 { return U32ToF64(BoolToU32(v)) }

// func S8ToF32(v int8) float32   { return S32ToF32(int32(v)) }
// func S8ToF64(v int8) float64   { return S32ToF64(int32(v)) }
// func S16ToF32(v int16) float32 { return S32ToF32(int32(v)) }
// func S16ToF64(v int16) float64 { return S32ToF64(int32(v)) }
func S32ToF32(v int32) float32 { return *(*float32)(unsafe.Pointer(&v)) }

// func S32ToF64(v int32) float64 { return float64(*(*float32)(unsafe.Pointer(&v))) } // FIXME: wrong
func S64ToF64(v int64) float64 { return *(*float64)(unsafe.Pointer(&v)) }

// func U8ToF32(v uint8) float32   { return U32ToF32(uint32(v)) }
// func U8ToF64(v uint8) float64   { return U32ToF64(uint32(v)) }
// func U16ToF32(v uint16) float32 { return U32ToF32(uint32(v)) }
// func U16ToF64(v uint16) float64 { return U32ToF64(uint32(v)) }
func U32ToF32(v uint32) float32 { return *(*float32)(unsafe.Pointer(&v)) }

// func U32ToF64(v uint32) float64 { return float64(*(*float32)(unsafe.Pointer(&v))) }
func U64ToF64(v uint64) float64 { return *(*float64)(unsafe.Pointer(&v)) }

func F32ToF64(v float32) float64 { return float64(v) }

func PointerToU32[T any](v *T) uint32 { return uint32(Uptr(v)) }

func Uptr[T any](v *T) uintptr { return uintptr(unsafe.Pointer(v)) }

// Lift functions
