package cm

// OKSizedResult represents a result sized to hold the OK type.
// The size of the OK type must be greater than or equal to the size of the Err type.
// For results with two zero-length types, use UnsizedResult.
type OKSizedResult[OK any, Err any] struct {
	Result[OK, OK, Err]
}

// ErrSizedResult represents a result sized to hold the Err type.
// The size of the Err type must be greater than or equal to the size of the OK type.
// For results with two zero-length types, use UnsizedResult.
type ErrSizedResult[OK any, Err any] struct {
	Result[Err, OK, Err]
}

// Result is a tagged union representing either the OK type or the Err type.
// Either OK or Err must have non-zero size, e.g. both cannot be struct{} or a zero-length array.
// For results with two zero-length types, use UnsizedResult.
// For results with no associated types, use UntypedResult.
type Result[Shape any, OK any, Err any] struct {
	v Variant2[Shape, OK, Err]
}

// IsErr returns true if r holds the error value.
func (r *Result[S, OK, Err]) IsErr() bool {
	return r.v.tag
}

// SetOK stores the OK value in r.
func (r *Result[S, OK, Err]) SetOK(ok OK) {
	r.v.Set0(ok)
}

// SetErr stores the error value in r.
func (r *Result[S, OK, Err]) SetErr(err Err) {
	r.v.Set1(err)
}

// OK returns the OK value for r and true if r represents the OK state.
// If r represents an error, then the zero value of OK is returned.
func (r *Result[S, OK, Err]) OK() (ok OK, isOK bool) {
	return r.v.Case0()
}

// Err returns the error value for r and true if r represents the error state.
// If r represents an OK value, then the zero value of Err is returned.
func (r *Result[S, OK, Err]) Err() (err Err, isErr bool) {
	return r.v.Case1()
}

// UnsizedResult is a tagged union representing either the OK type or the Err type.
// Both OK or Err must have zero size, e.g. both must be struct{} or a zero-length array.
// For results with one or more non-zero length types, use Result.
// For results with no associated types, use UntypedResult.
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
