package wjson

import (
	"encoding/json"
	"fmt"
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

// Ignore ignores the next logical JSON value, such as an object, array, or value.
func Ignore(dec *json.Decoder) error {
	return dec.Decode(ignore)
}

// DecodeObject decodes a single JSON object, passing each decoded key to f.
// Callback f is responsible for parsing the entire value of the object key
// before returning.
func DecodeObject(dec *json.Decoder, f func(key string) error) error {
	err := DecodeDelim(dec, '{')
	if err != nil {
		return err
	}
	for dec.More() {
		key, err := DecodeString(dec)
		if err != nil {
			return err
		}
		err = f(key)
		if err != nil {
			return err
		}
	}
	return DecodeDelim(dec, '}')
}

// DecodeArray decodes a single JSON array, passing the array index to f.
// Callback f is responsible for parsing the entire array element before returning.
func DecodeArray(dec *json.Decoder, f func(i int) error) error {
	err := DecodeDelim(dec, '[')
	if err != nil {
		return err
	}
	for i := 0; dec.More(); i++ {
		err = f(i)
		if err != nil {
			return err
		}
	}
	return DecodeDelim(dec, ']')
}

// DecodeDelim decodes a single JSON delimiter ({, }, [, ]). It will
// return an error if the next parsed token is not the expected delimiter.
func DecodeDelim(dec *json.Decoder, d json.Delim) error {
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
func DecodeString(dec *json.Decoder) (string, error) {
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
