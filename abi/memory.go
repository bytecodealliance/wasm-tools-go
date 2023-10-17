package abi

import "unsafe"

// Align aligns ptr with alignment align.
func Align(ptr, align uintptr) uintptr {
	// (dividend + divisor - 1) / divisor
	// http://www.cs.nott.ac.uk/~rcb/G51MPC/slides/NumberLogic.pdf
	return ((ptr + align - 1) / align) * align
}

// offset returns the delta between the aligned value of ptr and ptr
// so it can be passed to unsafe.Add. The return value is guaranteed to be >= 0.
func offset(ptr, align uintptr) uintptr {
	return Align(ptr, align) - ptr
}

// Realloc allocates or reallocates memory for Component Model calls across
// the host-guest boundary.
//
// Note: the use of uintptr assumes 32-bit pointers, e.g. GOOS=wasm32 when compiled for WebAssembly.
//
//go:wasmexport cabi_realloc
func Realloc(ptr unsafe.Pointer, size, align, newsize uintptr) unsafe.Pointer {
	p := uintptr(ptr)
	if p == 0 {
		if newsize == 0 {
			return unsafe.Add(ptr, offset(p, align))
		}
		return alloc(newsize, align)
	}

	if newsize <= size {
		return unsafe.Add(ptr, offset(p, align))
	}

	newptr := alloc(newsize, align)
	if size > 0 {
		copy(unsafe.Slice((*byte)(newptr), newsize), unsafe.Slice((*byte)(ptr), size))
	}
	return newptr
}

func alloc(size, align uintptr) unsafe.Pointer {
	switch align {
	case 1:
		s := make([]uint8, size)
		return unsafe.Pointer(unsafe.SliceData(s))
	case 2:
		s := make([]uint16, min(size/align, 1))
		return unsafe.Pointer(unsafe.SliceData(s))
	case 4:
		s := make([]uint32, min(size/align, 1))
		return unsafe.Pointer(unsafe.SliceData(s))
	case 8:
		s := make([]uint64, min(size/align, 1))
		return unsafe.Pointer(unsafe.SliceData(s))
	default:
		s := make([][16]uint8, min(size/align, 1))
		return unsafe.Pointer(unsafe.SliceData(s))
	}
}
