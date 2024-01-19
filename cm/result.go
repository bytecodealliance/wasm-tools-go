package cm

// Result is the common interface implemented by result types.
type Result[OK, Err any] interface {
	IsErr() bool
	SetOK(OK)
	SetErr(Err)
	OK() (ok OK, isOK bool)
	Err() (err Err, isErr bool)
}

// OKSizedResult represents a result sized to hold the OK type.
// The size of the OK type must be greater than or equal to the size of the Err type.
// For results with two zero-length types, use UnsizedResult.
type OKSizedResult[OK any, Err any] struct {
	SizedResult[OK, OK, Err]
}

// ErrSizedResult represents a result sized to hold the Err type.
// The size of the Err type must be greater than or equal to the size of the OK type.
// For results with two zero-length types, use UnsizedResult.
type ErrSizedResult[OK any, Err any] struct {
	SizedResult[Err, OK, Err]
}

// SizedResult is a tagged union that represents either the OK type or the Err type.
// Either OK or Err must have non-zero size, e.g. both cannot be struct{} or a zero-length array.
// For results with two zero-length types, use UnsizedResult.
type SizedResult[Shape any, OK any, Err any] struct {
	v Variant2[Shape, OK, Err]
}

// IsErr returns true if r holds the error value.
func (r *SizedResult[S, OK, Err]) IsErr() bool {
	return r.v.tag
}

// SetOK stores the OK value in r.
func (r *SizedResult[S, OK, Err]) SetOK(ok OK) {
	r.v.Set0(ok)
}

// SetErr stores the error value in r.
func (r *SizedResult[S, OK, Err]) SetErr(err Err) {
	r.v.Set1(err)
}

// OK returns the OK value for r and true if r represents the OK state.
// If r represents an error, then the zero value of OK is returned.
func (r *SizedResult[S, OK, Err]) OK() (ok OK, isOK bool) {
	return r.v.Case0()
}

// Err returns the error value for r and true if r represents the error state.
// If r represents an OK value, then the zero value of Err is returned.
func (r *SizedResult[S, OK, Err]) Err() (err Err, isErr bool) {
	return r.v.Case1()
}

type UnsizedResult[OK any, Err any] struct {
	v UnsizedVariant2[OK, Err]
}

// IsErr returns true if r holds the error value.
func (r *UnsizedResult[OK, Err]) IsErr() bool {
	return bool(r.v)
}

// SetErr stores the OK value in r.
func (r *UnsizedResult[OK, Err]) SetOK(ok OK) {
	r.v.Set0(ok)
}

// SetErr stores the error value in r.
func (r *UnsizedResult[OK, Err]) SetErr(err Err) {
	r.v.Set1(err)
}

// OK returns the OK value for r and true if r represents the OK state.
// If r represents an error, then the zero value of OK is returned.
func (r *UnsizedResult[OK, Err]) OK() (ok OK, isOK bool) {
	return r.v.Case0()
}

// Err returns the error value for r and true if r represents the error state.
// If r represents an OK value, then the zero value of Err is returned.
func (r *UnsizedResult[OK, Err]) Err() (err Err, isErr bool) {
	return r.v.Case1()
}

// UntypedResult represents an untyped Component Model result, e.g.
// result or result<_, _>. The OK and Err types are both struct{}.
type UntypedResult = UnsizedResult[struct{}, struct{}]
