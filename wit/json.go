package wit

import (
	"encoding/json"
	"fmt"
)

type unmarshalJSONFunc func(data []byte) error

func (f unmarshalJSONFunc) UnmarshalJSON(data []byte) error {
	return f(data)
}

func decodeIgnore(dec json.Decoder) error {
	//lint:ignore SA1014 this is a no-op
	return dec.Decode(unmarshalJSONFunc(func([]byte) error { return nil }))
}

func decodeObject(dec json.Decoder, f func(key string) error) error {
	err := decodeDelim(dec, '{')
	if err != nil {
		return err
	}
	for dec.More() {
		key, err := decodeString(dec)
		if err != nil {
			return err
		}
		err = f(key)
		if err != nil {
			return err
		}
	}
	return decodeDelim(dec, '}')
}

func decodeArray(dec json.Decoder, f func(i int) error) error {
	err := decodeDelim(dec, '[')
	if err != nil {
		return err
	}
	for i := 0; dec.More(); i++ {
		err = f(i)
		if err != nil {
			return err
		}
	}
	return decodeDelim(dec, ']')
}

func decodeDelim(dec json.Decoder, d json.Delim) error {
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

func decodeString(dec json.Decoder) (string, error) {
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
