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

type variant[Disc Discriminant, Shape, Align any] interface {
	isVariant(Disc, Shape, Align)
}

func (Variant[Disc, Shape, Align]) isVariant(Disc, Shape, Align) {}

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

// New returns a [Variant] with tag of type Disc, storage and GC shape of type Shape,
// aligned to type Align, with a value of type T.
func New[V variant[Disc, Shape, Align], Disc Discriminant, Shape, Align any, T any](tag Disc, data T) V {
	var v Variant[Disc, Shape, Align]
	if BoundsCheck && unsafe.Sizeof(data) > unsafe.Sizeof(v.data) {
		panic("NewVariant: size of requested type greater than size of data type")
	}
	v.tag = tag
	*(*T)(unsafe.Pointer(&v.data)) = data
	return *(*V)(unsafe.Pointer(&v))
}

// Tag returns the variant tag value.
func (v *Variant[Disc, Shape, Align]) Tag() Disc {
	return v.tag
}

// Set sets the [Variant] case to tag, and data to data.
func (v *Variant[Disc, Shape, Align]) Set(tag Disc, data unsafe.Pointer) {
	v.tag = tag
	v.data = *(*Shape)(data)
}

// Case returns a non-nil pointer if the [Variant] case is equal to tag, otherwise it returns nil.
func (v *Variant[Disc, Shape, Align]) Case(tag Disc) unsafe.Pointer {
	if v.tag == tag {
		return v.Data()
	}
	return nil
}

// Data returns an unsafe.Pointer to the data field in v.
func (v *Variant[Disc, Shape, Align]) Data() unsafe.Pointer {
	return unsafe.Pointer(&v.data)
}

// Case returns a non-nil *T if the [Variant] case is equal to tag, otherwise it returns nil.
func Case[T any, V variant[Disc, Shape, Align], Disc Discriminant, Shape, Align any](v *V, tag Disc) *T {
	if BoundsCheck {
		if unsafe.Sizeof(*(*T)(nil)) > unsafe.Sizeof(*(*Shape)(nil)) {
			panic("Get: size of requested type greater than size of data type")
		}
	}
	v2 := (*Variant[Disc, Shape, Align])(unsafe.Pointer(v))
	if v2.tag == tag {
		return (*T)(unsafe.Pointer(&v2.data))
	}
	return nil
}

// SetCase sets the [Variant] case to tag and data to data.
func SetCase[T any, V ~*Variant[Disc, Shape, Align], Disc Discriminant, Shape, Align any](v V, tag Disc, data T) {
	v2 := (*Variant[Disc, Shape, Align])(v)
	if BoundsCheck {
		if unsafe.Sizeof(*(*T)(nil)) > unsafe.Sizeof(v2.data) {
			panic("Get: size of requested type greater than size of data type")
		}
	}
	v2.tag = tag
	*((*T)(unsafe.Pointer(&v2.data))) = data
}
