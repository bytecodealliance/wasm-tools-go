package cm

import "unsafe"

const (
	// ResultOK represents the OK case of a result.
	ResultOK = false

	// ResultErr represents the error case of a result.
	ResultErr = true
)

// Result represents an untyped result, e.g. result or result<_, _>.
// Its associated types are implicitly struct{}, and it is represented as a Go bool.
type Result bool

// IsErr returns true if r represents the error case.
func (r Result) IsErr() bool {
	return bool(r)
}

// OK returns a non-nil pointer if r represents the OK case.
// If r represents an error, then it returns nil.
func (r Result) OK() *struct{} {
	if r {
		return nil
	}
	return &struct{}{}
}

// Err returns a non-nil pointer if r represents the error case.
// If r represents the OK case, then it returns nil.
func (r Result) Err() *struct{} {
	if !r {
		return nil
	}
	return &struct{}{}
}

// OKResult represents a result sized to hold the OK type.
// The size of the OK type must be greater than or equal to the size of the Err type.
// For results with two zero-length types, use [Result].
//
// TODO: change this to an alias when https://github.com/golang/go/issues/46477 is implemented.
type OKResult[OK, Err any] result[OK, OK, Err]

// IsErr returns true if r represents the error case.
func (r *OKResult[OK, Err]) IsErr() bool {
	return r.isErr
}

// OK returns a non-nil *OK pointer if r represents the OK case.
// If r represents an error, then it returns nil.
func (r *OKResult[OK, Err]) OK() *OK {
	if r.isErr {
		return nil
	}
	return (*OK)(unsafe.Pointer(&r.data))
}

// Err returns a non-nil *Err pointer if r represents the error case.
// If r represents the OK case, then it returns nil.
func (r *OKResult[OK, Err]) Err() *Err {
	if !r.isErr {
		return nil
	}
	return (*Err)(unsafe.Pointer(&r.data))
}

// ErrResult represents a result sized to hold the Err type.
// The size of the Err type must be greater than or equal to the size of the OK type.
// For results with two zero-length types, use [Result].
//
// TODO: change this to an alias when https://github.com/golang/go/issues/46477 is implemented.
type ErrResult[OK, Err any] result[Err, OK, Err]

// IsErr returns true if r represents the error case.
func (r *ErrResult[OK, Err]) IsErr() bool {
	return r.isErr
}

// OK returns a non-nil *OK pointer if r represents the OK case.
// If r represents an error, then it returns nil.
func (r *ErrResult[OK, Err]) OK() *OK {
	if r.isErr {
		return nil
	}
	return (*OK)(unsafe.Pointer(&r.data))
}

// Err returns a non-nil *Err pointer if r represents the error case.
// If r represents the OK case, then it returns nil.
func (r *ErrResult[OK, Err]) Err() *Err {
	if !r.isErr {
		return nil
	}
	return (*Err)(unsafe.Pointer(&r.data))
}

type result[Shape, OK, Err any] struct {
	isErr bool
	_     [0]OK
	_     [0]Err
	data  Shape
}

// OK returns an OK result with shape Shape and type OK and Err.
// Pass OKResult[OK, Err] or ErrResult[OK, Err] as the first type argument.
func OK[R ~struct {
	isErr bool
	_     [0]OK
	_     [0]Err
	data  Shape
}, Shape, OK, Err any](ok OK) R {
	if BoundsCheck && unsafe.Sizeof(*(*OK)(nil)) > unsafe.Sizeof(*(*Shape)(nil)) {
		panic("OK: size of requested type greater than size of data type")
	}
	var r result[Shape, OK, Err]
	r.isErr = ResultOK
	*((*OK)(unsafe.Pointer(&r.data))) = ok
	return R(r)
}

// Err returns an error result with shape Shape and type OK and Err.
// Pass OKResult[OK, Err] or ErrResult[OK, Err] as the first type argument.
func Err[R ~struct {
	isErr bool
	_     [0]OK
	_     [0]Err
	data  Shape
}, Shape, OK, Err any](err Err) R {
	if BoundsCheck && unsafe.Sizeof(*(*Err)(nil)) > unsafe.Sizeof(*(*Shape)(nil)) {
		panic("OK: size of requested type greater than size of data type")
	}
	var r result[Shape, OK, Err]
	r.isErr = ResultErr
	*((*Err)(unsafe.Pointer(&r.data))) = err
	return R(r)
}
