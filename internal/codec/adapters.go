package codec

// sliceCodec is an implementation of ElementDecoder for an arbitrary slice.
type sliceCodec[E comparable] []E

// Slice returns an ElementDecoder for slice s.
func Slice[E comparable](s *[]E) ElementDecoder {
	return (*sliceCodec[E])(s)
}

// DecodeElement implements the ElementDecoder interface,
// dynamically resizing the slice if necessary.
func (s *sliceCodec[E]) DecodeElement(dec Decoder, i int) error {
	var v E
	if i >= 0 && i < len(*s) {
		v = (*s)[i]
	}
	err := dec.Decode(&v)
	if err != nil {
		return err
	}
	Element(s, i)
	if v != (*s)[i] {
		(*s)[i] = v
	}
	return nil
}

// Element returns element i in slice s, reallocating the slice
// if necessary to at least len == i+1.
// TODO: rename to resize
func Element[S ~[]E, E comparable](s *S, i int) E {
	var e E
	if i < 0 {
		return e
	}
	if i >= len(*s) {
		*s = append(*s, make([]E, i+1-len(*s))...)
	}
	return (*s)[i]
}

// mapDecoder is an implementation of FieldDecoder for an arbitrary map with string keys.
type mapDecoder[K ~string, V any] map[K]V

// Map returns an FieldDecoder for map m.
func Map[K ~string, V any](m *map[K]V) FieldDecoder {
	return (*mapDecoder[K, V])(m)
}

// DecodeField implements the FieldDecoder interface,
// allocating the underlying map if necessary.
func (m *mapDecoder[K, V]) DecodeField(dec Decoder, name string) error {
	var v V
	err := dec.Decode(&v)
	if err != nil {
		return err
	}
	if *m == nil {
		*m = make(map[K]V)
	}
	(*m)[K(name)] = v
	return nil
}
