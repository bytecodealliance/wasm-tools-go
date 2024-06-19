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

// Lower functions
// func BoolToS32[B ~bool](v B) int32  { return int32(*(*int8)(unsafe.Pointer(&v))) }
func BoolToU32[B ~bool](v B) uint32 { return uint32(*(*uint8)(unsafe.Pointer(&v))) }
func LowerBool[B ~bool](v B) uint32 { return uint32(*(*uint8)(unsafe.Pointer(&v))) }

func BoolToU64[B ~bool](v B) uint64 { return uint64(*(*uint8)(unsafe.Pointer(&v))) }

func S32ToF32(v int32) float32   { return *(*float32)(unsafe.Pointer(&v)) }
func S64ToF64(v int64) float64   { return *(*float64)(unsafe.Pointer(&v)) }
func F64ToU64(v float64) uint64  { return *(*uint64)(unsafe.Pointer(&v)) }
func U32ToF32(v uint32) float32  { return *(*float32)(unsafe.Pointer(&v)) }
func U64ToF64(v uint64) float64  { return *(*float64)(unsafe.Pointer(&v)) }
func F32ToF64(v float32) float64 { return float64(v) }

func PointerToPointer[T any](v *T) *T { return v }
func PointerToU32[T any](v *T) uint32 { return uint32(uintptr(unsafe.Pointer(v))) }
func PointerToU64[T any](v *T) uint64 { return uint64(uintptr(unsafe.Pointer(v))) }

// Experimental lowering functions
func Lower1[T0, F0 CoreTypes, T any](v *T) F0 {
	p := (*struct {
		F0 T0
	})(unsafe.Pointer(v))
	return F0(p.F0)
}

func Lower2[T0, F0, T1, F1 CoreTypes, T any](v *T) (F0, F1) {
	p := (*struct {
		F0 T0
		F1 T1
	})(unsafe.Pointer(v))
	return F0(p.F0), F1(p.F1)
}
