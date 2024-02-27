package cm

import "unsafe"

const (
	// ResultOK represents the OK case of a result.
	ResultOK = false

	// ResultErr represents the error case of a result.
	ResultErr = true
)

// Result represents a result with no OK or error type.
// False represents the OK case and true represents the error case.
type Result bool

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
	return (*result[OK, OK, Err])(r).OK()
}

// Err returns a non-nil *Err pointer if r represents the error case.
// If r represents the OK case, then it returns nil.
func (r *OKResult[OK, Err]) Err() *Err {
	return (*result[OK, OK, Err])(r).Err()
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
	return (*result[Err, OK, Err])(r).OK()
}

// Err returns a non-nil *Err pointer if r represents the error case.
// If r represents the OK case, then it returns nil.
func (r *ErrResult[OK, Err]) Err() *Err {
	return (*result[Err, OK, Err])(r).Err()
}

type result[Shape, OK, Err any] struct {
	isErr bool
	_     [0]OK
	_     [0]Err
	data  Shape
}

func (r *result[Shape, OK, Err]) validate() {
	// Check if size of Shape is greater than both OK and Err
	if unsafe.Sizeof(*(*Shape)(nil)) > unsafe.Sizeof(*(*OK)(nil)) && unsafe.Sizeof(*(*Shape)(nil)) > unsafe.Sizeof(*(*Err)(nil)) {
		panic("result: size of data type > OK and Err types")
	}

	// Check if size of OK is greater than Shape
	if unsafe.Sizeof(*(*OK)(nil)) > unsafe.Sizeof(*(*Shape)(nil)) {
		panic("result: size of OK type > data type")
	}

	// Check if size of Err is greater than Shape
	if unsafe.Sizeof(*(*Err)(nil)) > unsafe.Sizeof(*(*Shape)(nil)) {
		panic("result: size of Err type > data type")
	}

	// Check if Shape is zero-sized, but align of OK is > 1
	if unsafe.Sizeof(*(*Shape)(nil)) == 0 && unsafe.Alignof(*(*OK)(nil)) != 1 {
		panic("result: size of data type == 0, but size != 1 (align of OK > 1)")
	}

	// Check if Shape is zero-sized, but align of Err is > 1
	if unsafe.Sizeof(*(*Shape)(nil)) == 0 && unsafe.Alignof(*(*Err)(nil)) != 1 {
		panic("result: size of data type == 0, but size != 1 (align of Err > 1)")
	}

	// Check if Shape is zero-sized, but size of result is > 1
	if unsafe.Sizeof(*(*Shape)(nil)) == 0 && unsafe.Sizeof(*r) != 1 {
		panic("result size != 1 (both OK and Err type struct{}?)")
	}
}

func (r *result[Shape, OK, Err]) IsErr() bool {
	r.validate()
	return r.isErr
}

func (r *result[Shape, OK, Err]) OK() *OK {
	r.validate()
	if r.isErr {
		return nil
	}
	return (*OK)(unsafe.Pointer(&r.data))
}

func (r *result[Shape, OK, Err]) Err() *Err {
	r.validate()
	if !r.isErr {
		return nil
	}
	return (*Err)(unsafe.Pointer(&r.data))
}

// OK returns an OK result with shape Shape and type OK and Err.
// Pass OKResult[OK, Err] or ErrResult[OK, Err] as the first type argument.
func OK[R ~struct {
	isErr bool
	_     [0]OK
	_     [0]Err
	data  Shape
}, Shape, OK, Err any](ok OK) R {
	var r result[Shape, OK, Err]
	r.validate()
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
	var r result[Shape, OK, Err]
	r.validate()
	r.isErr = ResultErr
	*((*Err)(unsafe.Pointer(&r.data))) = err
	return R(r)
}
