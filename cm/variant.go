package cm

import "unsafe"

// Discriminant is the set of types that can represent the tag or discriminator of a variant.
// Use bool for 2-value variants, results, or option<T> types, uint8 where there are 256 or
// fewer cases, uint16 for up to 65,536 cases, or uint32 for anything greater.s
type Discriminant interface {
	bool | uint8 | uint16 | uint32
}

// Variant represents a loosely-typed Component Model variant.
// Shape and Align must be non-zero sized types. To create a variant with no associated
// types, use UntypedVariant. To create a variant with only zero-sized types
// like struct{} or [0]int, use UnsizedVariant.
type Variant[Disc Discriminant, Shape, Align any] struct {
	tag  Disc
	_    [0]Align
	data Shape
}

// Tag returns the variant tag value.
func (v *Variant[Disc, Shape, Align]) Tag() Disc {
	return v.tag
}

// Data returns an unsafe.Pointer to the data field in v.
func (v *Variant[Disc, Shape, Align]) Data() unsafe.Pointer {
	return unsafe.Pointer(&v.data)
}

// Is returns true if the variant tag value equals tag
func (v *Variant[Disc, Shape, Align]) Is(tag Disc) bool {
	return v.tag == tag
}

// As coerces the value in Variant v into type T.
func As[T any, Disc Discriminant, Shape, Align any](v *Variant[Disc, Shape, Align]) (data T) {
	if BoundsCheck && unsafe.Sizeof(data) > unsafe.Sizeof(v.data) {
		panic("As: size of requested type greater than size of data type")
	}
	return *((*T)(unsafe.Pointer(&v.data)))
}

// Get returns a value of type T and true if the Variant tag is tag,
// or returns the zero value of T and false.
func Get[T any, Disc Discriminant, Shape, Align any](v *Variant[Disc, Shape, Align], tag Disc) (data T, ok bool) {
	if BoundsCheck && unsafe.Sizeof(data) > unsafe.Sizeof(v.data) {
		panic("Get: size of requested type greater than size of data type")
	}
	if v.tag != tag {
		var zero T
		return zero, false
	}
	return *((*T)(unsafe.Pointer(&v.data))), true
}

// Set sets the tag and value of Variant v.
func Set[T any, Disc Discriminant, Shape, Align any](v *Variant[Disc, Shape, Align], tag Disc, data T) {
	if BoundsCheck && unsafe.Sizeof(data) > unsafe.Sizeof(v.data) {
		panic("Set: size of requested type greater than size of data type")
	}
	v.tag = tag
	*((*T)(unsafe.Pointer(&v.data))) = data
}

// NewVariant returns a new variant with tag of type Disc, storage and GC shape of type Shape,
// setting the tag and value of type T.
func NewVariant[Disc Discriminant, Shape, Align any, T any](tag Disc, data T) Variant[Disc, Shape, Align] {
	var v Variant[Disc, Shape, Align]
	if BoundsCheck && unsafe.Sizeof(data) > unsafe.Sizeof(v.data) {
		panic("NewVariant: size of requested type greater than size of data type")
	}
	v.tag = tag
	*((*T)(unsafe.Pointer(&v.data))) = data
	return v
}

type Shape[T any] [1]T

// Variant2 represents a variant with 2 cases, where at least one case has an
// associated type with a non-zero size.
// Use UnsizedVariant2 if both T0 or T1 are zero-sized.
// The memory layout will have additional padding if both T0 and T1 are zero-sized.
type Variant2[Shape, T0, T1 any] struct {
	Variant[bool, Shape, align2[T0, T1]]
}

type align2[T0, T1 any] struct {
	_ [0]T0
	_ [0]T1
}

func (v *Variant2[S, T0, T1]) Case0() (val T0, ok bool) {
	return Get[T0](&v.Variant, false)
}

func (v *Variant2[S, T0, T1]) Case1() (val T1, ok bool) {
	return Get[T1](&v.Variant, true)
}

func (v *Variant2[S, T0, T1]) Set0(data T0) {
	Set[T0](&v.Variant, false, data)
}

func (v *Variant2[S, T0, T1]) Set1(data T1) {
	Set[T1](&v.Variant, true, data)
}

// UnsizedVariant2 represents a variant with 2 zero-sized associated types, e.g. struct{} or [0]T.
// Use Variant2 if either T0 or T1 has a non-zero size.
// Loads and stores may panic if T0 or T1 has a non-zero size.
type UnsizedVariant2[T0, T1 any] bool

func (v *UnsizedVariant2[T0, T1]) Tag() bool {
	return bool(*v)
}

func (v *UnsizedVariant2[T0, T1]) Case0() (data T0, ok bool) {
	return data, bool(*v)
}

func (v *UnsizedVariant2[T0, T1]) Case1() (data T1, ok bool) {
	return data, bool(*v)
}

func (v *UnsizedVariant2[T0, T1]) Set0(data T0) {
	*v = false
}

func (v *UnsizedVariant2[T0, T1]) Set1(data T1) {
	*v = true
}

// UntypedVariant2 represents an untyped variant with 2 cases.
// The associated types are defaulted to struct{}.
type UntypedVariant2 = UnsizedVariant2[struct{}, struct{}]
