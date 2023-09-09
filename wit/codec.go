package wit

import (
	"io"

	"github.com/ydnar/wit-bindgen-go/internal/codec/json"
)

func DecodeJSON(r io.Reader) (*Resolve, error) {
	res := &Resolve{}
	dec := json.NewDecoder(r, res)
	err := dec.Decode(res)
	return res, err
}

// Codec implements the codec.Codec interface,
// translating types to decoding/encoding-aware versions.
func (res *Resolve) Codec(v any) (any, error) {
	switch v := v.(type) {
	// WIT sections
	case *[]*World:
		return (*ptrSliceCodec[World])(v), nil
	case *[]*Interface:
		return (*ptrSliceCodec[Interface])(v), nil
	case *[]*TypeDef:
		return (*ptrSliceCodec[TypeDef])(v), nil
	case *[]*Package:
		return (*ptrSliceCodec[Package])(v), nil

	// Maps of concrete types
	case map[string]*Function:
		return mapCodec[Function](v), nil

	// Maps of references
	case map[string]WorldItem:
		return mapFuncCodec[WorldItem](v), nil
	case map[string]*Interface:
		return mapFuncCodec[*Interface](v), nil
	case map[string]*World:
		return mapFuncCodec[*World](v), nil

	// Callback setters of enum/interface types
	case func(WorldItem):
		return worldItemCodec(v), nil
	case func(Type):
		return &typeCodec{v, res}, nil

	// Callback setters of references
	case func(*World):
		return &indexCodec[World]{&res.Worlds, v}, nil
	case func(*Interface):
		return &indexCodec[Interface]{&res.Interfaces, v}, nil
	case func(*TypeDef):
		return &indexCodec[TypeDef]{&res.TypeDefs, v}, nil
	case func(*Package):
		return &indexCodec[Package]{&res.Packages, v}, nil
	}
	return nil, nil
}

func (c *Resolve) DecodeField(name string) (any, error) {
	switch name {
	case "worlds":
		return &c.Worlds, nil
	case "interfaces":
		return &c.Interfaces, nil
	case "types":
		return &c.TypeDefs, nil
	case "packages":
		return &c.Packages, nil
	}
	return nil, nil
}

func (c *World) DecodeField(name string) (any, error) {
	switch name {
	case "name":
		return &c.Name, nil
	case "docs":
		return &c.Docs, nil
	case "imports":
		return &c.Imports, nil
	case "exports":
		return &c.Exports, nil
	}
	return nil, nil
}

type worldItemCodec func(WorldItem)

func (c worldItemCodec) DecodeField(name string) (any, error) {
	switch name {
	case "interface":
		return func(i *Interface) { c(i) }, nil
	case "type":
		return func(i *TypeDef) { c(i) }, nil
	}
	return nil, nil
}

func (i *Interface) DecodeField(name string) (any, error) {
	switch name {
	case "docs":
		return &i.Docs, nil
	case "name":
		return &i.Name, nil
	case "types":
		return &i.Types, nil
	case "functions":
		return &i.Functions, nil
	case "package":
		return &i.Package, nil
	}
	return nil, nil
}

// typeCodec translates WIT type strings or reference IDs into a Type.
// The caller is responsible for setting f to receive the resolved Type.
type typeCodec struct {
	f func(Type)
	*Resolve
}

// DecodeString translates a into to a primitive WIT type.
// c.f is called with the resulting Type, if any.
func (c typeCodec) DecodeString(s string) error {
	t, err := ParseType(s)
	if err != nil {
		return err
	}
	c.f(t)
	return nil
}

// DecodeInt translates a TypeDef reference into a pointer to a TypeDef
// in the parent Resolve struct.
func (c typeCodec) DecodeInt(i int64) error {
	var t Type = realloc(&c.TypeDefs, int(i))
	c.f(t)
	return nil
}

func (f *Function) DecodeElement(name string) (any, error) {
	switch name {
	case "docs":
		return &f.Docs, nil
	case "name":
		return &f.Name, nil
	case "kind":
		return &f.Kind, nil
	case "params":
		return &f.Params, nil
	case "results":
		return &f.Results, nil
	}
	return nil, nil
}

type indexCodec[T any] struct {
	s *[]*T
	f func(*T)
}

func (c *indexCodec[T]) DecodeInt(i int64) error {
	v := realloc(c.s, int(i))
	c.f(v)
	return nil
}

type ptrSliceCodec[T any] []*T

func (c *ptrSliceCodec[T]) DecodeElement(i int) (any, error) {
	return realloc(c, i), nil
}

type sliceCodec[T any] []T

func (c *sliceCodec[T]) DecodeElement(i int) (any, error) {
	remake(c, i)
	return &(*c)[i], nil
}

// remake reallocates a slice if necessary for len i,
// returning the value of s[i].
func remake[S ~[]E, E any](s *S, i int) E {
	var e E
	if i < 0 {
		return e
	}
	if i >= len(*s) {
		*s = append(*s, make([]E, i+1-len(*s))...)
	}
	return (*s)[i]
}

// realloc returns the value of slice s at index i,
// reallocating the slice if necessary. s must be a slice
// of pointers, because the underlying backing to s might
// change when reallocated.
// If the value at s[i] is nil, a new *E will be allocated.
func realloc[S ~[]*E, E any](s *S, i int) *E {
	if i < 0 {
		return nil
	}
	remake(s, i)
	if (*s)[i] == nil {
		(*s)[i] = new(E)
	}
	return (*s)[i]
}

type mapCodec[T any] map[string]*T

func (c *mapCodec[T]) DecodeField(name string) (any, error) {
	if _, ok := (*c)[name]; !ok {
		if *c == nil {
			*c = make(map[string]*T)
		}
		(*c)[name] = new(T)
	}
	return (*c)[name], nil
}

type mapFuncCodec[T any] map[string]T

func (c *mapFuncCodec[T]) DecodeField(name string) (any, error) {
	return func(v T) {
		if *c == nil {
			*c = make(map[string]T)
		}
		(*c)[name] = v
	}, nil
}
