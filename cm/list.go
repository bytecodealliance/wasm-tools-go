package cm

import "unsafe"

// List represents a Component Model list<T>.
// The binary representation of list<T> is similar to a Go slice minus the cap field.
type List[T any] struct {
	ptr *T
	len uintptr
}

// ToList returns a List[T] equivalent to the Go slice s.
// The underlying slice data is not copied, and the resulting List points at the
// same array storage as the slice.
func ToList[S ~[]T, T any](s S) List[T] {
	return List[T]{
		ptr: unsafe.SliceData([]T(s)),
		len: uintptr(len(s)),
	}
}

// Slice returns a Go slice representing the List.
func (list List[T]) Slice() []T {
	return unsafe.Slice(list.ptr, list.len)
}
