package codec

type Visitor[C any] interface {
	Visitor(ctx C) (any, error)
}

// Visit returns the visitor for v with context C.
// If v does not implement Visitor, v is returned unchanged.
func Visit[C any](ctx C, v any) (any, error) {
	if v, ok := v.(Visitor[C]); ok {
		return v.Visitor(ctx)
	}
	return v, nil
}

type Decoder interface {
	Decode(v any) error
}

type NilDecoder interface {
	DecodeNil() error
}

type Value interface {
	bool | int64 | uint64 | float64 | string | []byte
}

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
