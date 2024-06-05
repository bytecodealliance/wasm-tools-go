package cm

import (
	"fmt"
	"unsafe"
)

type CoreTypes interface {
	int32 | uint32 | int64 | uint64 | float32 | float64
}

func Lower[To CoreTypes, From any](from From) To {
	var to To
	to -= 1
	fmt.Printf("%x\t", uint64(to))

	switch unsafe.Sizeof(to) {
	case 4:
		if to > 0 {
			fmt.Printf("%x\n", *(*uint32)(unsafe.Pointer(&to)))
		} else {
			fmt.Printf("%x\n", *(*int32)(unsafe.Pointer(&to)))
		}
	case 8:
		if to > 0 {
			fmt.Printf("%x\n", *(*uint64)(unsafe.Pointer(&to)))
		} else {
			fmt.Printf("%x\n", *(*int64)(unsafe.Pointer(&to)))
		}
	}
	switch (*uint32)(unsafe.Pointer(&to)) {
	}

	// switch unsafe.Sizeof(from) {
	// case 1:
	// 	return To((*uint8)(unsafe.Pointer(&from)))
	// }

	return to
}

func StringData(s string) uintptr {
	return uintptr(unsafe.Pointer(unsafe.StringData(s)))
}

func AnyToU32[T any](v T) uint32 { return *(*uint32)(unsafe.Pointer(&v)) }

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
