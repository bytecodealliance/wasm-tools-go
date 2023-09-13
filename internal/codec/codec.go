package codec

// Resolver is the interface implemented by types that return a codec for the value at v.
// Values returned by Resolver should implement one or more encode or decode methods.
type Resolver interface {
	ResolveCodec(v any) any
}

// Resolvers is a slice of Resolver values. It also implements the Resolver interface.
type Resolvers []Resolver

// ResolveCodec walks the slice of Resolvers, returning the first non-nil value or an error.
func (rs Resolvers) ResolveCodec(v any) any {
	for _, r := range rs {
		c := r.ResolveCodec(v)
		if c != nil {
			return c
		}
	}
	return nil
}

// Decoder is the interface implemented by types that can decode data into Go type(s).
type Decoder interface {
	Decode(v any) error
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

type IntDecoder[T Integer] interface {
	DecodeInt(T) error
}

type FloatDecoder[T Float] interface {
	DecodeFloat(T) error
}

type FieldDecoder interface {
	DecodeField(dec Decoder, name string) error
}

type ElementDecoder interface {
	DecodeElement(dec Decoder, i int) error
}
