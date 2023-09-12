package codec

import "fmt"

// Resolver is the interface implemented by types that return a codec for the value at v.
// Values returned by Resolver should implement one or more encode or decode methods.
type Resolver interface {
	ResolveCodec(v any) (any, error)
}

// Resolvers is a slice of Resolver values. It also implements the Resolver interface.
type Resolvers []Resolver

// ResolveCodec walks the slice of Resolvers, returning the first non-nil value or an error.
func (rs Resolvers) ResolveCodec(v any) (any, error) {
	for _, r := range rs {
		c, err := r.ResolveCodec(v)
		if err != nil {
			return nil, err
		}
		if c != nil {
			return c, nil
		}
	}
	return nil, nil
}

// Decoder is the interface implemented by types that can decode data into Go type(s).
type Decoder interface {
	Decode(v any) error
}

// EndDecoder is the interface implemented by types that wish to receive a signal
// that decoding has finished. DecodeEnd is not called if an error occurs during
// decoding. DecodeEnd can return an error to abort further decoding.
type EndDecoder interface {
	DecodeEnd() error
}

type NilDecoder interface {
	DecodeNil() error
}

type BoolDecoder interface {
	DecodeBool(bool) error
}

type BytesDecoder interface {
	DecodeBytes([]byte) error
}

type StringDecoder interface {
	DecodeString(string) error
}

type IntDecoder interface {
	DecodeInt(int64) error
}

type UintDecoder interface {
	DecodeUint(uint64) error
}

type FloatDecoder interface {
	DecodeFloat(float64) error
}

type FieldDecoder interface {
	DecodeField(dec Decoder, name string) error
}
type ElementDecoder interface {
	DecodeElement(dec Decoder, i int) error
}

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
	Index(s, i)
	if v != (*s)[i] {
		(*s)[i] = v
	}
	return nil
}

// Index returns element i in slice s, reallocating the slice
// if necessary to at least len == i+1.
func Index[S ~[]E, E any](s *S, i int) E {
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
	fmt.Printf("(*Map).DecodeField(dec, %q) -> %#v\n", name, v)
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
