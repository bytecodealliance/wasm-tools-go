package wit

// Align aligns ptr with alignment align.
func Align(ptr, align uintptr) uintptr {
	// (dividend + divisor - 1) / divisor
	// http://www.cs.nott.ac.uk/~rcb/G51MPC/slides/NumberLogic.pdf
	return ((ptr + align - 1) / align) * align
}

// Discriminant returns the smallest integer type that can represent 0...n.
func Discriminant(n int) Type {
	switch {
	case n <= 1<<8:
		return U8{}
	case n <= 1<<16:
		return U16{}
	}
	return U32{}
}

// Sized is the interface implemented by any type that reports its [ABI byte size] and [alignment].
//
// [ABI byte size]: https://github.com/WebAssembly/component-model/blob/main/design/mvp/CanonicalABI.md#size
// [alignment]: https://github.com/WebAssembly/component-model/blob/main/design/mvp/CanonicalABI.md#alignment
type Sized interface {
	Size() uintptr
	Align() uintptr
}

type _sized struct{}

func (_sized) Size() uintptr  { panic("unimplemented") }
func (_sized) Align() uintptr { panic("unimplemented") }
