package cm

import "unsafe"

const (
	// ResultOK represents the OK case of a result.
	ResultOK = false

	// ResultErr represents the error case of a result.
	ResultErr = true
)

// OKResult represents a result with only an OK type.
//
// TODO: change this to an alias when https://github.com/golang/go/issues/46477 is implemented.
type OKResult[OK any] struct {
	Result[OK, OK, struct{}]
}

// ErrResult represents a result with only an error type.
//
// TODO: change this to an alias when https://github.com/golang/go/issues/46477 is implemented.
type ErrResult[Err any] struct {
	Result[Err, struct{}, Err]
}

// OKSizedResult represents a result sized to hold the OK type.
// The size of the OK type must be greater than or equal to the size of the Err type.
// For results with two zero-length types, use UnsizedResult.
//
// TODO: change this to an alias when https://github.com/golang/go/issues/46477 is implemented.
type OKSizedResult[OK any, Err any] struct {
	Result[OK, OK, Err]
}

// ErrSizedResult represents a result sized to hold the Err type.
// The size of the Err type must be greater than or equal to the size of the OK type.
// For results with two zero-length types, use UnsizedResult.
//
// TODO: change this to an alias when https://github.com/golang/go/issues/46477 is implemented.
type ErrSizedResult[OK any, Err any] struct {
	Result[Err, OK, Err]
}

// Result is a tagged union representing either the OK type or the Err type.
// Either OK or Err must have non-zero size, e.g. both cannot be struct{} or a zero-length array.
// For results with zero-sized or no associated types, use [UntypedResult].
type Result[Shape, OK, Err any] struct {
	isErr bool
	_     [0]OK
	_     [0]Err
	data  Shape
}

// OK returns an OK [Result] with GC shape Shape and type OK and Err.
func OK[Shape, OK, Err any](ok OK) Result[Shape, OK, Err] {
	var r Result[Shape, OK, Err]
	if BoundsCheck && unsafe.Sizeof(ok) > unsafe.Sizeof(r.data) {
		panic("OK: size of requested type greater than size of data type")
	}
	r.isErr = ResultOK
	*((*OK)(unsafe.Pointer(&r.data))) = ok
	return r
}

// Err returns an error [Result] with GC shape Shape and type OK and Err.
func Err[Shape, OK, Err any](err Err) Result[Shape, OK, Err] {
	var r Result[Shape, OK, Err]
	if BoundsCheck && unsafe.Sizeof(err) > unsafe.Sizeof(r.data) {
		panic("Err: size of requested type greater than size of data type")
	}
	r.isErr = ResultErr
	*((*Err)(unsafe.Pointer(&r.data))) = err
	return r
}

// IsErr returns true if r represents the error case.
func (r *Result[Shape, OK, Err]) IsErr() bool {
	return r.isErr
}

// OK returns a non-nil *OK pointer if r represents the OK case.
// If r represents an error, then it returns nil.
func (r *Result[Shape, OK, Err]) OK() *OK {
	if r.isErr {
		return nil
	}
	return (*OK)(unsafe.Pointer(&r.data))
}

// Err returns a non-nil *Err pointer if r represents the error case.
// If r represents the OK case, then it returns nil.
func (r *Result[Shape, OK, Err]) Err() *Err {
	if !r.isErr {
		return nil
	}
	return (*Err)(unsafe.Pointer(&r.data))
}

// UntypedResult represents an untyped result, e.g. result or result<_, _>.
// Its associated types are implicitly struct{}, and it is represented as a Go bool.
type UntypedResult bool

// IsErr returns true if r represents the error case.
func (r UntypedResult) IsErr() bool {
	return bool(r)
}

// OK returns a non-nil pointer if r represents the OK case.
// If r represents an error, then it returns nil.
func (r UntypedResult) OK() *struct{} {
	if r {
		return nil
	}
	return &struct{}{}
}

// Err returns a non-nil pointer if r represents the error case.
// If r represents the OK case, then it returns nil.
func (r UntypedResult) Err() *struct{} {
	if !r {
		return nil
	}
	return &struct{}{}
}
