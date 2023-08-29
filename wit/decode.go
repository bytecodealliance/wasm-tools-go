package wit

import (
	"encoding/json"
	"io"

	"github.com/ydnar/wit-bindgen-go/internal/wjson"
)

func DecodeJSON(r io.Reader) (*Resolve, error) {
	dec := &decodeState{
		dec: json.NewDecoder(r),
		res: &Resolve{},
	}
	dec.dec.UseNumber()

	err := dec.decodeResolve()
	if err != nil {
		return nil, err
	}

	return dec.res, nil
}

type decodeState struct {
	dec *json.Decoder
	res *Resolve
}

func (dec *decodeState) decodeResolve() error {
	return wjson.DecodeObject(dec.dec, func(key string) error {
		switch key {
		case "worlds":
			return dec.decodeResolveItem(func(i int) error {
				return dec.decodeWorld(element(&dec.res.Worlds, i))
			})
		case "interfaces":
			return dec.decodeResolveItem(func(i int) error {
				return dec.decodeInterface(element(&dec.res.Interfaces, i))
			})
		case "types":
			return dec.decodeResolveItem(func(i int) error {
				return dec.decodeTypeDef(element(&dec.res.Types, i))
			})
		case "packages":
			return dec.decodeResolveItem(func(i int) error {
				return dec.decodePackage(element(&dec.res.Packages, i))
			})
		default:
			return wjson.Ignore(dec.dec)
		}
	})
}

func (dec *decodeState) decodeResolveItem(f func(i int) error) error {
	return wjson.DecodeObject(dec.dec, func(key string) error {
		switch key {
		case "items":
			return wjson.DecodeArray(dec.dec, f)
		default:
			return wjson.Ignore(dec.dec)
		}
	})
}

func (dec *decodeState) decodeWorld(world *World) error {
	return wjson.Ignore(dec.dec)
}

func (dec *decodeState) decodeInterface(iface *Interface) error {
	return wjson.Ignore(dec.dec)
}

func (dec *decodeState) decodeTypeDef(typ *TypeDef) error {
	return wjson.Ignore(dec.dec)
}

func (dec *decodeState) decodePackage(pkg *Package) error {
	return wjson.Ignore(dec.dec)
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
