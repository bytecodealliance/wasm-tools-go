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
	var v Variant[Disc, Shape, Align]
	if BoundsCheck && unsafe.Sizeof(data) > unsafe.Sizeof(v.data) {
		panic("NewVariant: size of requested type greater than size of data type")
	}
	v.tag = tag
	*((*T)(unsafe.Pointer(&v.data))) = data
	return v
}

// Tag returns the variant tag value.
func (v *Variant[Disc, Shape, Align]) Tag() Disc {
	return v.tag
}

// Data returns an unsafe.Pointer to the data field in v.
func (v *Variant[Disc, Shape, Align]) Data() unsafe.Pointer {
	return unsafe.Pointer(&v.data)
}

// Case returns a non-nil *T if the [Variant] case is equal to tag, otherwise it returns nil.
func Case[T any, V ~*Variant[Disc, Shape, Align], Disc Discriminant, Shape, Align any](v V, tag Disc) *T {
	v1 := (*Variant[Disc, Shape, Align])(v)
	if BoundsCheck {
		if unsafe.Sizeof(*(*T)(nil)) > unsafe.Sizeof(v1.data) {
			panic("Get: size of requested type greater than size of data type")
		}
	}
	if v1.tag == tag {
		return (*T)(unsafe.Pointer(&v1.data))
	}
	return nil
}
