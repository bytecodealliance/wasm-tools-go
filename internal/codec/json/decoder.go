package json

import (
	"encoding"
	"encoding/json"
	"fmt"
	"io"

	"github.com/ydnar/wit-bindgen-go/internal/codec"
	"golang.org/x/exp/constraints"
)

type Decoder struct {
	dec *json.Decoder
}

var _ codec.Decoder = &Decoder{}

func NewDecoder(r io.Reader) *Decoder {
	dec := json.NewDecoder(r)
	dec.UseNumber()
	return &Decoder{
		dec: dec,
	}
}

func (dec *Decoder) Decode(v any) error {
	tok, err := dec.dec.Token()
	if err == io.EOF {
		return nil
	}
	if err != nil {
		return err
	}

	if tok == nil {
		return dec.decodeNull(v)
	}

	switch tok := tok.(type) {
	case bool:
		return dec.decodeBool(v, tok)
	case json.Number:
		return dec.decodeNumber(v, tok)
	case string:
		return dec.decodeString(v, tok)
	case json.Delim:
		switch tok {
		case '{':
			return dec.decodeObject(v)
		case '[':
			return dec.decodeArray(v)
		default:
			return fmt.Errorf("unexpected JSON token %v at offset %d", tok, dec.dec.InputOffset())
		}
	}

	return nil
}

func (dec *Decoder) decodeNull(v any) error {
	if v, ok := v.(codec.NilDecoder); ok {
		return v.DecodeNil()
	}
	return nil
}

func (dec *Decoder) decodeBool(v any, b bool) error {
	if v, ok := v.(codec.ValueDecoder[bool]); ok {
		return v.DecodeValue(b)
	}
	return nil
}

func (dec *Decoder) decodeNumber(v any, n json.Number) error {
	switch v := v.(type) {
	case *int64:
		return coerceNumber(v, n.Int64)
	case *uint64:
		return coerceNumber(v, n.Int64)
	case *float64:
		return coerceNumber(v, n.Float64)
	case *string:
		*v = string(n)
		// TODO: how to handle undecodable types?
		// Return an error? Silently ignore? Configurable?
	}
	return nil
}

func coerceNumber[To, From constraints.Integer | constraints.Float](v *To, f func() (From, error)) error {
	n, err := f()
	if err != nil {
		return err
	}
	*v = To(n)
	return nil
}

func (dec *Decoder) decodeString(v any, s string) error {
	switch v := v.(type) {
	case *string:
		*v = s
	case codec.ValueDecoder[string]:
		return v.DecodeValue(s)
	case encoding.TextUnmarshaler:
		return v.UnmarshalText([]byte(s))
		// TODO: how to handle undecodable types?
		// Return an error? Silently ignore? Configurable?
	}
	return nil
}

// decodeObject decodes a JSON object into v.
// It expects that the initial { token has already been decoded.
func (dec *Decoder) decodeObject(o any) error {
	d, ok := o.(codec.FieldDecoder)
	if !ok {
		// TODO: how to handle undecodable objects?
		d = &ignore{}
	}

	for dec.dec.More() {
		name, err := dec.stringToken()
		if err != nil {
			return err
		}
		v, err := d.DecodeField(name)
		if err != nil {
			return err
		}
		err = dec.Decode(v)
		if err != nil {
			return err
		}
	}
	return nil
}

// decodeArray decodes a JSON array into v.
// It expects that the initial [ token has already been decoded.
func (dec *Decoder) decodeArray(v any) error {
	d, ok := v.(codec.ElementDecoder)
	if !ok {
		// TODO: how to handle undecodable arrays?
		d = &ignore{}
	}

	for i := 0; dec.dec.More(); i++ {
		v, err := d.DecodeElement(i)
		if err != nil {
			return err
		}
		err = dec.Decode(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (dec *Decoder) stringToken() (string, error) {
	tok, err := dec.dec.Token()
	if err != nil {
		return "", err
	}
	s, ok := tok.(string)
	if !ok {
		return "", fmt.Errorf("unexpected JSON token %v at offset %d", tok, dec.dec.InputOffset())
	}
	return s, nil
}

type ignore struct{}

func (i *ignore) DecodeField(string) (any, error) {
	return i, nil
}

func (i *ignore) DecodeElement(int) (any, error) {
	return i, nil
}

func (i *ignore) UnmarshalJSON([]byte) error {
	return nil
}
