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
	// References
	case **World:
		return &worldCodec{v, res}, nil
	case **Interface:
		return &interfaceCodec{v, res}, nil
	case **TypeDef:
		return &typeDefCodec{v, res}, nil
	case **Package:
		return &packageCodec{v, res}, nil

	// Handles
	case **Function:
		newIfNil(v)
		return *v, nil

	// Enums
	case *Type:
		return &typeCodec{v, res}, nil
	case *TypeOwner:
		return &typeOwnerCodec{v}, nil
	case *WorldItem:
		return &worldItemCodec{v}, nil
	}

	return nil, nil
}

func (c *Resolve) getWorld(i int) *World {
	return newIfElementNil(&c.Worlds, i)
}

func (c *Resolve) getInterface(i int) *Interface {
	return newIfElementNil(&c.Interfaces, i)
}

func (c *Resolve) getTypeDef(i int) *TypeDef {
	return newIfElementNil(&c.TypeDefs, i)
}

func (c *Resolve) getPackage(i int) *Package {
	return newIfElementNil(&c.Packages, i)
}

func (c *Resolve) DecodeField(dec codec.Decoder, name string) error {
	switch name {
	case "worlds":
		return codec.DecodeSlice(dec, &c.Worlds)
	case "interfaces":
		return codec.DecodeSlice(dec, &c.Interfaces)
	case "types":
		return codec.DecodeSlice(dec, &c.TypeDefs)
	case "packages":
		return codec.DecodeSlice(dec, &c.Packages)
	}
	return nil
}

type worldCodec struct {
	w **World
	*Resolve
}

func (c *worldCodec) DecodeInt(i int) error {
	*c.w = c.getWorld(i)
	return nil
}

func (c *worldCodec) DecodeField(dec codec.Decoder, name string) error {
	w := newIfNil(c.w)
	switch name {
	case "name":
		return dec.Decode(&w.Name)
	case "docs":
		return dec.Decode(&w.Docs)
	case "imports":
		return codec.DecodeMap(dec, &w.Imports)
	case "exports":
		return codec.DecodeMap(dec, &w.Exports)
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

type interfaceCodec struct {
	i **Interface
	*Resolve
}

func (c *interfaceCodec) DecodeInt(i int) error {
	*c.i = c.getInterface(i)
	return nil
}

func (c *interfaceCodec) DecodeField(dec codec.Decoder, name string) error {
	i := newIfNil(c.i)
	switch name {
	case "docs":
		return dec.Decode(&i.Docs)
	case "name":
		return dec.Decode(&i.Name)
	case "types":
		return codec.DecodeMap(dec, &i.TypeDefs)
	case "functions":
		return codec.DecodeMap(dec, &i.Functions)
	case "package":
		return dec.Decode(&i.Package)
	}
	return nil
}

type typeDefCodec struct {
	t **TypeDef
	*Resolve
}

func (c *typeDefCodec) DecodeInt(i int) error {
	*c.t = c.getTypeDef(i)
	return nil
}

func (c *typeDefCodec) DecodeField(dec codec.Decoder, name string) error {
	t := newIfNil(c.t)
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
	var err error
	*c.t, err = ParseType(s)
	return err
}

// DecodeInt translates a TypeDef reference into a pointer to a TypeDef
// in the parent Resolve struct.
func (c *typeCodec) DecodeInt(i int) error {
	*c.t = c.getTypeDef(i)
	return nil
}

// typeOwnerCodec translates WIT type owner enums into a TypeOwner.
type typeOwnerCodec struct {
	o *TypeOwner
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
		return codec.DecodeSlice(dec, &f.Params)
	case "results":
		return codec.DecodeSlice(dec, &f.Results)
	}
	return nil
}

type packageCodec struct {
	p **Package
	*Resolve
}

func (c *packageCodec) DecodeInt(i int) error {
	*c.p = c.getPackage(i)
	return nil
}

func (c *packageCodec) DecodeField(dec codec.Decoder, name string) error {
	p := newIfNil(c.p)
	switch name {
	case "docs":
		return dec.Decode(&p.Docs)
	case "name":
		return dec.Decode(&p.Name)
	case "interfaces":
		return codec.DecodeMap(dec, &p.Interfaces)
	case "worlds":
		return codec.DecodeMap(dec, &p.Worlds)
	}
	return nil
}

// newIfNil allocates a new instance of T if *v == nil.
func newIfNil[T any](v **T) *T {
	if *v == nil {
		*v = new(T)
	}
	return *v
}

// newIfElementNil resizes s and allocates a new instance of T if necessary.
func newIfElementNil[S ~[]*E, E any](s *S, i int) *E {
	if i < 0 {
		return nil
	}
	if codec.Resize(s, i) == nil {
		(*s)[i] = new(E)
	}
	return (*s)[i]
}
