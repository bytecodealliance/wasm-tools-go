package abi

import "unsafe"

// Align aligns ptr with alignment align.
func Align(ptr unsafe.Pointer, align uintptr) unsafe.Pointer {
	// (dividend + divisor - 1) / divisor
	// http://www.cs.nott.ac.uk/~rcb/G51MPC/slides/NumberLogic.pdf
	p := (uintptr(unsafe.Add(ptr, align-1)) / align) * align
	// Appease vet, see https://github.com/golang/go/issues/58625
	return *(*unsafe.Pointer)(unsafe.Pointer(&p))
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
			return unsafe.Pointer(Align(ptr, align))
		}
		return alloc(newsize, align)
	}

	if newsize <= size {
		return unsafe.Pointer(Align(ptr, align))
	}

	newptr := alloc(newsize, align)
	if size > 0 {
		// Appease vet, see https://github.com/golang/go/issues/58625
		src := *(*unsafe.Pointer)(unsafe.Pointer(&ptr))
		copy(unsafe.Slice((*byte)(newptr), newsize), unsafe.Slice((*byte)(src), size))
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
