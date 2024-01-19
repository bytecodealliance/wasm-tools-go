package cm

import "unsafe"

type Discriminant interface {
	uint8 | uint16 | uint32
}

// Variant represents a loosely-typed Component Model variant.
// Disc must be one of uint8, uint16, or uint32. Shape and Align
// must be non-zero sized types. To create a variant with no associated
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

func init() {
	var v Variant[uint8, uint64, uint64]
	Set(&v, 1, uint64(99))
	_, _ = Get[int64](&v, 1)
	_ = NewVariant[uint8, string, string](0, "hello world")
	_ = NewVariant[uint8, string, string](0, "hello world")
}

type Shape[T any] [1]T

type Variant2[T0, T1 any] interface {
	V() uint
	V0() (T0, bool)
	V1() (T1, bool)
	Set0(T0)
	Set1(T1)
}

// SizedVariant2 represents a variant with 2 associated types, at least one of which has a non-zero size.
// Use UnsizedVariant2 if both T0 or T1 are zero-sized.
// The memory layout will have additional padding if both T0 and T1 are zero-sized.
type SizedVariant2[S Shape[T0] | Shape[T1], T0, T1 any] struct {
	isT1 bool
	_    [0]T0
	_    [0]T1
	val  S
}

func (v *SizedVariant2[S, T0, T1]) V() uint {
	return uint(*(*uint8)(unsafe.Pointer(&v.isT1)))
}

func (v *SizedVariant2[S, T0, T1]) V0() (val T0, ok bool) {
	return *(*T0)(unsafe.Pointer(&v.val)), !v.isT1

}

func (v *SizedVariant2[S, T0, T1]) V1() (val T1, ok bool) {
	return *(*T1)(unsafe.Pointer(&v.val)), v.isT1
}

func (v *SizedVariant2[S, T0, T1]) Set0(val T0) {
	v.isT1 = false
	*(*T0)(unsafe.Pointer(&v.val)) = val
}

func (v *SizedVariant2[S, T0, T1]) Set1(val T1) {
	v.isT1 = true
	*(*T1)(unsafe.Pointer(&v.val)) = val
}

// UnsizedVariant2 represents a variant with 2 zero-sized associated types, e.g. struct{} or [0]T.
// Use SizedVariant2 if either T0 or T1 has a non-zero size.
// Loads and stores may panic if T0 or T1 has a non-zero size.
type UnsizedVariant2[T0, T1 any] bool

func (v *UnsizedVariant2[T0, T1]) V() uint {
	return uint(*(*uint8)(unsafe.Pointer(v)))
}

func (v *UnsizedVariant2[T0, T1]) V0() (val T0, ok bool) {
	return val, bool(*v)
}

func (v *UnsizedVariant2[T0, T1]) V1() (val T1, ok bool) {
	return val, bool(*v)
}

func (v *UnsizedVariant2[T0, T1]) Set0(val T0) {
	*v = false
}

func (v *UnsizedVariant2[T0, T1]) Set1(val T1) {
	*v = true
}

// UntypedVariant2 represents an untyped variant of cardinality 2.
// The associated types are defaulted to struct{}.
type UntypedVariant2 = UnsizedVariant2[struct{}, struct{}]
