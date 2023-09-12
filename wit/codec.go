package wit

import (
	"io"

	"github.com/ydnar/wit-bindgen-go/internal/codec"
	"github.com/ydnar/wit-bindgen-go/internal/codec/json"
)

func DecodeJSON(r io.Reader) (*Resolve, error) {
	res := &Resolve{}
	dec := json.NewDecoder(r, res)
	err := dec.Decode(res)
	return res, err
}

// ResolveCodec implements the codec.Resolver interface,
// translating types to decoding/encoding-aware versions.
func (res *Resolve) ResolveCodec(v any) (any, error) {
	switch v := v.(type) {
	// WIT sections
	case *[]*World:
		return codec.AsSlice(v), nil
	case *[]*Interface:
		return codec.AsSlice(v), nil
	case *[]*TypeDef:
		return codec.AsSlice(v), nil
	case *[]*Package:
		return codec.AsSlice(v), nil

	// Maps
	case *map[string]WorldItem:
		return codec.AsMap(v), nil
	case *map[string]*Function:
		return codec.AsMap(v), nil
	case *map[string]*Interface:
		return codec.AsMap(v), nil
	case *map[string]*World:
		return codec.AsMap(v), nil

	// References
	case **World:
		return asRefCodec(v, &res.Worlds), nil
	case **Interface:
		return asRefCodec(v, &res.Interfaces), nil
	case **TypeDef:
		return asRefCodec(v, &res.TypeDefs), nil
	case **Package:
		return asRefCodec(v, &res.Packages), nil

	// Handles
	case **Function:
		if *v == nil {
			*v = new(Function)
		}
		return *v, nil

	// Enums
	case *Type:
		return &typeCodec{v, res}, nil
	case *TypeOwner:
		return &typeOwnerCodec{v, res}, nil
	case *WorldItem:
		return &worldItemCodec{v}, nil
	}

	return nil, nil
}

func (c *Resolve) DecodeField(dec codec.Decoder, name string) error {
	switch name {
	case "worlds":
		return dec.Decode(&c.Worlds)
	case "interfaces":
		return dec.Decode(&c.Interfaces)
	case "types":
		return dec.Decode(&c.TypeDefs)
	case "packages":
		return dec.Decode(&c.Packages)
	}
	return nil
}

func (w *World) DecodeField(dec codec.Decoder, name string) error {
	switch name {
	case "name":
		return dec.Decode(&w.Name)
	case "docs":
		return dec.Decode(&w.Docs)
	case "imports":
		return dec.Decode(&w.Imports)
	case "exports":
		return dec.Decode(&w.Exports)
	}
	return nil
}

type worldItemCodec struct {
	i *WorldItem
}

func (c worldItemCodec) DecodeField(dec codec.Decoder, name string) error {
	switch name {
	case "interface":
		var i *Interface
		err := dec.Decode(&i)
		if err != nil {
			return err
		}
		*c.i = i
	case "type":
		var t *TypeDef
		err := dec.Decode(&t)
		if err != nil {
			return err
		}
		*c.i = t
	}
	return nil
}

func (i *Interface) DecodeField(dec codec.Decoder, name string) error {
	switch name {
	case "docs":
		return dec.Decode(&i.Docs)
	case "name":
		return dec.Decode(&i.Name)
	case "types":
		return dec.Decode(&i.TypeDefs)
	case "functions":
		return dec.Decode(&i.Functions)
	case "package":
		return dec.Decode(&i.Package)
	}
	return nil
}

func (t *TypeDef) DecodeField(dec codec.Decoder, name string) error {
	switch name {
	case "kind":
		return dec.Decode(&t.Kind)
	case "name":
		return dec.Decode(&t.Name)
	case "owner":
		return dec.Decode(&t.Owner)
	}
	return nil
}

// typeCodec translates WIT type strings or reference IDs into a Type.
type typeCodec struct {
	t *Type
	*Resolve
}

// DecodeString translates a into to a primitive WIT type.
// c.f is called with the resulting Type, if any.
func (c *typeCodec) DecodeString(s string) error {
	t, err := ParseType(s)
	if err != nil {
		return err
	}
	*c.t = t
	return nil
}

// DecodeInt translates a TypeDef reference into a pointer to a TypeDef
// in the parent Resolve struct.
func (c *typeCodec) DecodeInt(i int64) error {
	t := codec.Element(&c.TypeDefs, int(i))
	if t == nil {
		t = new(TypeDef)
		c.TypeDefs[i] = t
	}
	*c.t = t
	return nil
}

// typeOwnerCodec translates WIT type owner enums into a TypeOwner.
type typeOwnerCodec struct {
	o *TypeOwner
	*Resolve
}

func (c *typeOwnerCodec) DecodeField(dec codec.Decoder, name string) error {
	switch name {
	case "interface":
		var i *Interface
		err := dec.Decode(&i)
		if err != nil {
			return err
		}
		*c.o = i
	case "world":
		var w *World
		err := dec.Decode(&w)
		if err != nil {
			return err
		}
		*c.o = w
	}
	return nil
}

func (f *Function) DecodeField(dec codec.Decoder, name string) error {
	switch name {
	case "docs":
		return dec.Decode(&f.Docs)
	case "name":
		return dec.Decode(&f.Name)
	case "kind":
		return dec.Decode(&f.Kind)
	case "params":
		return dec.Decode(&f.Params)
	case "results":
		return dec.Decode(&f.Results)
	}
	return nil
}

func (p *Package) DecodeField(dec codec.Decoder, name string) error {
	switch name {
	case "docs":
		return dec.Decode(&p.Docs)
	case "name":
		return dec.Decode(&p.Name)
	case "interfaces":
		return dec.Decode(&p.Interfaces)
	case "worlds":
		return dec.Decode(&p.Worlds)
	}
	return nil
}

// refCodec is a codec for decoding a 0-based reference to
// some resource (T), or data structure representing that
// type.
type refCodec[T any] struct {
	v **T
	s *[]*T
}

func asRefCodec[T any](v **T, s *[]*T) *refCodec[T] {
	return &refCodec[T]{v, s}
}

func (c *refCodec[T]) DecodeInt(i int64) error {
	*c.v = codec.Element(c.s, int(i))
	if *c.v == nil {
		*c.v = new(T)
		(*c.s)[i] = *c.v
	}
	return nil
}

func (c *refCodec[T]) DecodeElement(dec codec.Decoder, i int) error {
	if *c.v == nil {
		*c.v = new(T)
	}
	if ed, ok := any(*c.v).(codec.ElementDecoder); ok {
		return ed.DecodeElement(dec, i)
	}
	return nil
}

func (c *refCodec[T]) DecodeField(dec codec.Decoder, name string) error {
	if *c.v == nil {
		*c.v = new(T)
	}
	if fd, ok := any(*c.v).(codec.FieldDecoder); ok {
		return fd.DecodeField(dec, name)
	}
	return nil
}
