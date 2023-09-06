package wit

import (
	"encoding/json"
	"io"

	"github.com/ydnar/wit-bindgen-go/internal/wjson"
)

func (res *Resolve) DecodeField(name string) (any, error) {
	switch name {
	case "worlds":
		return &res.Worlds, nil
	}
	return nil, nil
}

type worldDecoder struct {
	*World
	res *Resolve
}

func (w *worldDecoder) DecodeField(name string) (any, error) {
	switch name {
	case "name":
		return &w.Name, nil
	case "docs":
		return &w.Docs, nil
	case "imports":
		w.Imports = make(map[string]WorldItem)
		return &worldItemsDecoder{w.Imports, w.res}, nil
	case "exports":
		w.Exports = make(map[string]WorldItem)
		return &worldItemsDecoder{w.Exports, w.res}, nil
	}
	return nil, nil
}

type worldItemsDecoder struct {
	m   map[string]WorldItem
	res *Resolve
}

func (m *worldItemsDecoder) DecodeField(name string) (any, error) {
	return nil, nil
}

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
