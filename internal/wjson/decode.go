package wjson

import (
	"encoding/json"
	"fmt"

	"golang.org/x/exp/constraints"
)

// UnmarshalFunc is a function that implements the json.Unmarshaler interface.
type UnmarshalFunc func(data []byte) error

// UnmarshalJSON implements json.Unmarshaler.
func (f UnmarshalFunc) UnmarshalJSON(data []byte) error {
	if f == nil {
		return nil
	}
	return f(data)
}

var ignore json.Unmarshaler = UnmarshalFunc(nil)

// Decoder is any type that implements this subset of json.Decoder methods.
type Decoder interface {
	Decode(any) error
	More() bool
	Token() (json.Token, error)
	InputOffset() int64
}

// Ignore ignores the next logical JSON value, such as an object, array, or value.
func Ignore(dec Decoder) error {
	return dec.Decode(ignore)
}

// DecodeObject decodes a single JSON object, passing each decoded key to f.
// Callback f must parse nothing, or parse the entire field.
// Partially parsing a field, or parsing more than the field is undefined.
// If f parses nothing, the field will be ignored.
func DecodeObject(dec Decoder, f func(key string) error) error {
	err := DecodeDelim(dec, '{')
	if err != nil {
		return err
	}
	for dec.More() {
		key, err := DecodeString(dec)
		if err != nil {
			return err
		}
		offset := dec.InputOffset()
		err = f(key)
		if err != nil {
			return err
		}
		if dec.InputOffset() == offset {
			Ignore(dec)
		}
	}
	return DecodeDelim(dec, '}')
}

// DecodeArray decodes a single JSON array, passing the array index to f.
// Callback f must parse nothing, or parse the entire array element.
// Partially parsing an element, or parsing more than the element is undefined.
// If f parses nothing, the element will be ignored.
func DecodeArray(dec Decoder, f func(i int) error) error {
	err := DecodeDelim(dec, '[')
	if err != nil {
		return err
	}
	for i := 0; dec.More(); i++ {
		offset := dec.InputOffset()
		err = f(i)
		if err != nil {
			return err
		}
		if dec.InputOffset() == offset {
			Ignore(dec)
		}
	}
	return DecodeDelim(dec, ']')
}

// DecodeDelim decodes a single JSON delimiter ({, }, [, ]). It will
// return an error if the next parsed token is not the expected delimiter.
func DecodeDelim(dec Decoder, d json.Delim) error {
	tok, err := dec.Token()
	if err != nil {
		return err
	}
	got, ok := tok.(json.Delim)
	if !ok || got != d {
		return fmt.Errorf("expected %v, got %v", d, tok)
	}
	return nil
}

// DecodeString decodes a single JSON string. It will return an error if
// the next parsed token is not a string.
func DecodeString(dec Decoder) (string, error) {
	tok, err := dec.Token()
	if err != nil {
		return "", err
	}
	s, ok := tok.(string)
	if !ok {
		return "", fmt.Errorf("expected string, got %v", tok)
	}
	return s, nil
}

// DecodeOptionalString decodes a single JSON string. It will return an error if
// the next parsed token is not a string or null.
func DecodeOptionalString(dec Decoder) (*string, error) {
	tok, err := dec.Token()
	if err != nil {
		return nil, err
	}
	if tok == nil {
		return nil, nil
	}
	s, ok := tok.(string)
	if !ok {
		return nil, fmt.Errorf("expected string, got %v", tok)
	}
	return &s, nil
}

// DecodeNumber decodes a single JSON number. It will return an error if
// the next parsed token is not a number or null.
func DecodeNumber[T constraints.Integer | constraints.Float](dec Decoder) (*T, error) {
	tok, err := dec.Token()
	if err != nil {
		return nil, err
	}
	if tok == nil {
		return nil, nil
	}
	var v T
	switch tok := tok.(type) {
	case float64:
		v = T(tok)
	case json.Number:
		switch any(v).(type) {
		case float32, float64: // TODO: this doesn’t cover ~float32 or ~float64
			f, err := tok.Float64()
			if err != nil {
				return nil, err
			}
			v = T(f)
		default:
			f, err := tok.Int64()
			if err != nil {
				return nil, err
			}
			v = T(f)
		}

	default:
		return nil, fmt.Errorf("expected number, got %v", tok)
	}
	return &v, nil
}

// DecodeInt decodes a single JSON integer. It will return an error if
// the next parsed token is not a number.
func DecodeInt[T constraints.Integer](dec Decoder, v *T) error {
	tok, err := dec.Token()
	if err != nil {
		return err
	}
	if tok == nil {
		return nil
	}
	switch tok := tok.(type) {
	case float64:
		*v = T(tok)
	case json.Number:
		i, err := tok.Int64()
		if err != nil {
			return err
		}
		*v = T(i)
	default:
		return fmt.Errorf("expected number, got %v", tok)
	}
	return nil
}
