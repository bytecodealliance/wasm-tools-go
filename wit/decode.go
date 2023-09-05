package wit

import (
	"encoding/json"
	"io"

	"github.com/ydnar/wit-bindgen-go/internal/wjson"
)

func DecodeJSON(r io.Reader) (*Resolve, error) {
	dec := &decodeState{
		Decoder: json.NewDecoder(r),
		res:     &Resolve{},
	}
	dec.UseNumber()

	err := dec.decodeResolve()
	if err != nil {
		return nil, err
	}

	return dec.res, nil
}

type decodeState struct {
	*json.Decoder
	res *Resolve
}

func (dec *decodeState) decodeResolve() error {
	return wjson.DecodeObject(dec, func(key string) error {
		switch key {
		case "worlds":
			return wjson.DecodeArray(dec, func(i int) error {
				return dec.decodeWorld(remake(&dec.res.Worlds, i))
			})
		case "interfaces":
			return wjson.DecodeArray(dec, func(i int) error {
				return dec.decodeInterface(remake(&dec.res.Interfaces, i))
			})
		case "types":
			return wjson.DecodeArray(dec, func(i int) error {
				return dec.decodeTypeDef(remake(&dec.res.TypeDefs, i))
			})
		case "packages":
			return wjson.DecodeArray(dec, func(i int) error {
				return dec.decodePackage(remake(&dec.res.Packages, i))
			})
		default:
			return nil
		}
	})
}

func (dec *decodeState) decodeWorld(world *World) error {
	return wjson.DecodeObject(dec, func(key string) error {
		switch key {
		case "name":
			return dec.Decode(&world.Name)
		case "docs":
			return dec.Decode(&world.Docs)
		case "package":
			return decodeIndex(dec, &dec.res.Packages, &world.Package)
		default:
			return nil
		}
	})
}

func (dec *decodeState) decodeInterface(iface *Interface) error {
	return nil
}

func (dec *decodeState) decodeTypeDef(typ *TypeDef) error {
	return nil
}

func (dec *decodeState) decodePackage(pkg *Package) error {
	return nil
}

func decodeIndex[S ~[]*E, E any](dec *decodeState, s *S, e **E) error {
	var i int
	err := wjson.DecodeInt(dec, &i)
	if err != nil {
		return err
	}
	*e = remake(s, i)
	return nil
}

// remake returns the value of slice s at index i,
// reallocating the slice if necessary. s must be a slice
// of pointers, because the underlying backing to s might
// change when reallocated.
// If the value at s[i] is nil, a new *E will be allocated.
func remake[S ~[]*E, E any](s *S, i int) *E {
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
