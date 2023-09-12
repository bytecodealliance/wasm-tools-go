package codec

// Slice is an implementation of ElementDecoder for an arbitrary slice.
type Slice[E comparable] []E

// AsSlice returns s coerced into a Slice.
func AsSlice[E comparable](s *[]E) *Slice[E] {
	return (*Slice[E])(s)
}

// DecodeElement implements the ElementDecoder interface,
// dynamically resizing the slice if necessary.
func (s *Slice[E]) DecodeElement(dec Decoder, i int) error {
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

// Map is an implementation of FieldDecoder for an arbitrary map with string keys.
type Map[K ~string, V any] map[K]V

// AsMap returns m coerced into a Map.
func AsMap[K ~string, V any](m *map[K]V) *Map[K, V] {
	return (*Map[K, V])(m)
}

// DecodeField implements the FieldDecoder interface,
// allocating the underlying map if necessary.
func (m *Map[K, V]) DecodeField(dec Decoder, name string) error {
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

type IntDecoderFunc[T Integer] func(T) error

func (f IntDecoderFunc[T]) DecodeInt(v T) error {
	return f(v)
}

type FloatDecoderFunc[T Float] func(T) error

func (f FloatDecoderFunc[T]) DecodeFloat(v T) error {
	return f(v)
}
