package codec

// Resolver is the interface implemented by types that return a codec for the value at v.
// Values returned by Resolver should implement one or more encode or decode methods.
type Resolver interface {
	Resolve(v any) (any, error)
}

// Resolvers is a slice of Resolver values. It also implements the Resolver interface.
type Resolvers []Resolver

// Resolve walks the list of Resolvers, returning the first non-nil value or an error.
func (rs Resolvers) Resolve(v any) (any, error) {
	for _, r := range rs {
		c, err := r.Resolve(v)
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
