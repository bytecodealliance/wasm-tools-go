package wit

import (
	"encoding/json"
	"io"

	"github.com/ydnar/wit-bindgen-go/internal/wjson"
)

func DecodeJSON(r io.Reader) (*Resolve, error) {
	res := &Resolve{}

	dec := json.NewDecoder(r)
	dec.UseNumber()

	err := decodeResolve(dec, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func decodeResolve(dec *json.Decoder, res *Resolve) error {
	return wjson.DecodeObject(dec, func(key string) error {
		switch key {
		case "worlds":
			return decodeResolveItem(dec, res, func(i int) error {
				return decodeWorld(dec, res, element(&res.Worlds, i))
			})
		case "interfaces":
			return decodeResolveItem(dec, res, func(i int) error {
				return decodeInterface(dec, res, element(&res.Interfaces, i))
			})
		case "types":
			return decodeResolveItem(dec, res, func(i int) error {
				return decodeTypeDef(dec, res, element(&res.Types, i))
			})
		case "packages":
			return decodeResolveItem(dec, res, func(i int) error {
				return decodePackage(dec, res, element(&res.Packages, i))
			})
		default:
			return wjson.Ignore(dec)
		}
	})
}

func decodeResolveItem(dec *json.Decoder, res *Resolve, f func(i int) error) error {
	return wjson.DecodeObject(dec, func(key string) error {
		switch key {
		case "items":
			return wjson.DecodeArray(dec, f)
		default:
			return wjson.Ignore(dec)
		}
	})
}

func decodeWorld(dec *json.Decoder, res *Resolve, world *World) error {
	return wjson.Ignore(dec)
}

func decodeInterface(dec *json.Decoder, res *Resolve, iface *Interface) error {
	return wjson.Ignore(dec)
}

func decodeTypeDef(dec *json.Decoder, res *Resolve, typ *TypeDef) error {
	return wjson.Ignore(dec)
}

func decodePackage(dec *json.Decoder, res *Resolve, pkg *Package) error {
	return wjson.Ignore(dec)
}

// element returns the value of slice s at index i,
// reallocating the slice if necessary. s must be a slice
// of pointers, because the underlying backing to s might
// change when reallocated.
// If the value at s[i] is nil, a new *E will be allocated.
func element[S ~[]*E, E any](s *S, i int) *E {
	if i < 0 {
		return nil
	}
	if i >= len(*s) {
		*s = append(*s, make([]*E, i-len(*s))...)
	}
	if (*s)[i] == nil {
		(*s)[i] = new(E)
	}
	return (*s)[i]
}
