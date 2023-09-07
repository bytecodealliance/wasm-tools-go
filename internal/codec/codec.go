package codec

// Codec is the interface implemented by types that return a codec for the value at v.
// Values returned by Codec should implement one or more encoder or decoder methods.
type Codec interface {
	Codec(v any) (any, error)
}

type NilDecoder interface {
	DecodeNil() error
}

// TODO: delete this
type Value interface {
	bool | int64 | uint64 | float64 | string | []byte
}

// TODO: StringDecoder, IntDecoder, FloatDecoder
type ValueDecoder[T Value] interface {
	DecodeValue(v T) error
}

type FieldDecoder interface {
	DecodeField(name string) (any, error)
}

type FieldDecoderFunc func(name string) (any, error)

func (f FieldDecoderFunc) DecodeField(name string) (any, error) {
	return f(name)
}

type ElementDecoder interface {
	DecodeElement(i int) (any, error)
}

type ElementDecoderFunc func(i int) (any, error)

func (f ElementDecoderFunc) DecodeElement(i int) (any, error) {
	return f(i)
}
