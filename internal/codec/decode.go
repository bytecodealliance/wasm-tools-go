package codec

import (
	"encoding"
	"strconv"
)

// DecodeNil calls DecodeNil on v if v implements NilDecoder.
func DecodeNil(v any) error {
	if v, ok := v.(NilDecoder); ok {
		return v.DecodeNil()
	}
	return nil
}

// DecodeBool decodes a boolean value into v.
// If v implements BoolDecoder, then DecodeBool(b) is called.
func DecodeBool(v any, b bool) error {
	switch v := v.(type) {
	case *bool:
		*v = b
	case BoolDecoder:
		return v.DecodeBool(b)
	}
	return nil
}

// DecodeNumber decodes a number encoded as a string into v.
// The following types are supported: int64, uint64, float64, IntDecoder, and FloatDecoder.
// It will also attempt to decode the string directly into v if it is one of:
// string, []byte, StringDecoder, BytesDecoder, or encoding.TextUnmarshaler.
func DecodeNumber(v any, n string) error {
	switch v := v.(type) {
	case *int64:
		i, err := strconv.ParseInt(n, 10, 64)
		if err != nil {
			return err
		}
		*v = i
	case *uint64:
		i, err := strconv.ParseUint(n, 10, 64)
		if err != nil {
			return err
		}
		*v = i
	case *float64:
		f, err := strconv.ParseFloat(n, 64)
		if err != nil {
			return err
		}
		*v = f
	case IntDecoder:
		i, err := strconv.ParseInt(n, 10, 64)
		if err != nil {
			return err
		}
		return v.DecodeInt(i)
	case UintDecoder:
		i, err := strconv.ParseUint(n, 10, 64)
		if err != nil {
			return err
		}
		return v.DecodeUint(i)
	case FloatDecoder:
		f, err := strconv.ParseFloat(n, 64)
		if err != nil {
			return err
		}
		return v.DecodeFloat(f)
	}
	// TODO: how to handle undecodable types?
	// Return an error? Silently ignore? Configurable?
	return DecodeString(v, n)
}

// Decode string decodes s into v. The following types are supported:
// string, []byte, StringDecoder, BytesDecoder, and encoding.TextUnmarshaler.
func DecodeString(v any, s string) error {
	switch v := v.(type) {
	case *string:
		*v = s
	case *[]byte:
		*v = []byte(s)
	case StringDecoder:
		return v.DecodeString(s)
	case BytesDecoder:
		return v.DecodeBytes([]byte(s))
	case encoding.TextUnmarshaler:
		return v.UnmarshalText([]byte(s))
		// TODO: how to handle undecodable types?
		// Return an error? Silently ignore? Configurable?
	}
	return nil
}
