package codec

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

// type MapDecoder[T any] interface {
// 	Decode(m map[string]T) error
// }

type T struct {
	Name   string
	Height float64
	Weight float64
}

func (t *T) DecodeField(name string) (any, error) {
	switch name {
	case "name":
		return &t.Name, nil
	case "height":
		return &t.Height, nil
	case "weight":
		return &t.Weight, nil
	}
	return nil, nil
}
