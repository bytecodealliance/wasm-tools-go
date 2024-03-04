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
// types, use an enum.
type Variant[Disc Discriminant, Shape, Align any] struct {
	tag  Disc
	_    [0]Align
	data Shape
}

// This function is sized so it can be inlined and optimized away.
func validate[Disc Discriminant, Shape, Align any, T any]() {
	var v Variant[Disc, Shape, Align]
	var t T

	// Check if size of T is greater than Shape
	if unsafe.Sizeof(t) > unsafe.Sizeof(v.data) {
		panic("result: size of requested type > data type")
	}

	// Check if Shape is zero-sized, but size of result != 1
	if unsafe.Sizeof(v.data) == 0 && unsafe.Sizeof(v) != 1 {
		panic("result: size of data type == 0, but result size != 1")
	}
}

// NewVariant returns a [Variant] with tag of type Disc, storage and GC shape of type Shape,
// aligned to type Align, with a value of type T.
func NewVariant[Disc Discriminant, Shape, Align any, T any](tag Disc, data T) Variant[Disc, Shape, Align] {
	validate[Disc, Shape, Align, T]()
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
	validate[Disc, Shape, Align, T]()
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
	validate[Disc, Shape, Align, struct{}]()
	v2 := (*Variant[Disc, Shape, Align])(unsafe.Pointer(v))
	return v2.tag
}

// Is returns true if the [Variant] case is equal to tag.
func Is[V ~struct {
	tag  Disc
	_    [0]Align
	data Shape
}, Disc Discriminant, Shape, Align any](v *V, tag Disc) bool {
	validate[Disc, Shape, Align, struct{}]()
	return (*Variant[Disc, Shape, Align])(unsafe.Pointer(v)).tag == tag
}

// Case returns a non-nil *T if the [Variant] case is equal to tag, otherwise it returns nil.
func Case[T any, V ~struct {
	tag  Disc
	_    [0]Align
	data Shape
}, Disc Discriminant, Shape, Align any](v *V, tag Disc) *T {
	validate[Disc, Shape, Align, T]()
	v2 := (*Variant[Disc, Shape, Align])(unsafe.Pointer(v))
	if v2.tag == tag {
		return (*T)(unsafe.Pointer(&v2.data))
	}
	return nil
}
