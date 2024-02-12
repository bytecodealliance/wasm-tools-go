package cm

import "unsafe"

// Discriminant is the set of types that can represent the tag or discriminator of a variant.
// Use bool for 2-value variants, results, or option<T> types, uint8 where there are 256 or
// fewer cases, uint16 for up to 65,536 cases, or uint32 for anything greater.
type Discriminant interface {
	bool | uint8 | uint16 | uint32
}

// Variant represents a loosely-typed Component Model variant.
// Shape and Align must be non-zero sized types. To create a variant with no associated
// types, use UntypedVariant. For variants with zero-width or no associated types, use an enum.
type Variant[Disc Discriminant, Shape, Align any] struct {
	tag  Disc
	_    [0]Align
	data Shape
}

// NewVariant returns a [Variant] with tag of type Disc, storage and GC shape of type Shape,
// aligned to type Align, with a value of type T.
func NewVariant[Disc Discriminant, Shape, Align any, T any](tag Disc, data T) Variant[Disc, Shape, Align] {
	if BoundsCheck && unsafe.Sizeof(*(*T)(nil)) > unsafe.Sizeof(*(*Shape)(nil)) {
		panic("NewVariant: size of requested type greater than size of data type")
	}
	var v Variant[Disc, Shape, Align]
	v.tag = tag
	v.data = *(*Shape)(unsafe.Pointer(&data))
	return v
}

// New returns a [Variant] with tag of type Disc, storage and GC shape of type Shape,
// aligned to type Align, with a value of type T.
func New[V ~struct {
	tag  Disc
	_    [0]Align
	data Shape
}, Disc Discriminant, Shape, Align any, T any](tag Disc, data T) V {
	if BoundsCheck && unsafe.Sizeof(*(*T)(nil)) > unsafe.Sizeof(*(*Shape)(nil)) {
		panic("NewVariant: size of requested type greater than size of data type")
	}
	var v Variant[Disc, Shape, Align]
	v.tag = tag
	v.data = *(*Shape)(unsafe.Pointer(&data))
	return V(v)
}

// Tag returns the tag of [Variant] v.
func Tag[V ~struct {
	tag  Disc
	_    [0]Align
	data Shape
}, Disc Discriminant, Shape, Align any](v *V) Disc {
	v2 := (*Variant[Disc, Shape, Align])(unsafe.Pointer(v))
	return v2.tag
}

// Is returns true if the [Variant] case is equal to tag.
func Is[V ~struct {
	tag  Disc
	_    [0]Align
	data Shape
}, Disc Discriminant, Shape, Align any](v *V, tag Disc) bool {
	return (*Variant[Disc, Shape, Align])(unsafe.Pointer(v)).tag == tag
}

// Case returns a non-nil *T if the [Variant] case is equal to tag, otherwise it returns nil.
func Case[T any, V ~struct {
	tag  Disc
	_    [0]Align
	data Shape
}, Disc Discriminant, Shape, Align any](v *V, tag Disc) *T {
	if BoundsCheck && unsafe.Sizeof(*(*T)(nil)) > unsafe.Sizeof(*(*Shape)(nil)) {
		panic("Get: size of requested type greater than size of data type")
	}
	v2 := (*Variant[Disc, Shape, Align])(unsafe.Pointer(v))
	if v2.tag == tag {
		return (*T)(unsafe.Pointer(&v2.data))
	}
	return nil
}
