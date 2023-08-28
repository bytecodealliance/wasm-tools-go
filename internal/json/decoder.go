package json

import (
	"encoding/json"
	"fmt"
	"io"
)

// UnmarshalJSONFunc is a function that implements the json.Unmarshaler interface.
type UnmarshalJSONFunc func(data []byte) error

// UnmarshalJSON implements json.Unmarshaler.
func (f UnmarshalJSONFunc) UnmarshalJSON(data []byte) error {
	return f(data)
}

var _ json.Unmarshaler = UnmarshalJSONFunc(nil)

// Decoder wraps json.Decoder with additional features.
type Decoder json.Decoder

// NewDecoder returns a new Decoder that reads from r.
func NewDecoder(r io.Reader) *Decoder {
	return (*Decoder)(json.NewDecoder(r))
}

func (dec *Decoder) Ignore() error {
	//lint:ignore SA1014 this is a no-op
	return dec.Decode(UnmarshalJSONFunc(func([]byte) error { return nil }))
}

func (dec *Decoder) Decode(v any) error {
	return (*json.Decoder)(dec).Decode(v)
}

func (dec *Decoder) More() bool {
	return (*json.Decoder)(dec).More()
}

// DecodeObject decodes a single JSON object, passing each decoded key to f.
// Callback f is responsible for parsing the entire value of the object key
// before returning.
func (dec *Decoder) DecodeObject(f func(key string) error) error {
	err := dec.DecodeDelim('{')
	if err != nil {
		return err
	}
	for dec.More() {
		key, err := dec.DecodeString()
		if err != nil {
			return err
		}
		err = f(key)
		if err != nil {
			return err
		}
	}
	return dec.DecodeDelim('}')
}

// DecodeArray decodes a single JSON array, passing the array index to f.
// Callback f is responsible for parsing the entire array element before returning.
func (dec *Decoder) DecodeArray(f func(i int) error) error {
	err := dec.DecodeDelim('[')
	if err != nil {
		return err
	}
	for i := 0; dec.More(); i++ {
		err = f(i)
		if err != nil {
			return err
		}
	}
	return dec.DecodeDelim(']')
}

// DecodeDelim decodes a single JSON delimiter ({, }, [, ]). It will
// return an error if the next parsed token is not the expected delimiter.
func (dec *Decoder) DecodeDelim(d json.Delim) error {
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
// the next parsed token is not a string,
func (dec *Decoder) DecodeString() (string, error) {
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
