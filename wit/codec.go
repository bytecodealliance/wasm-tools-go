package wit

type Codec Resolve

func (c *Codec) Codec(v any) (any, error) {
	switch v := v.(type) {
	// Outer Resolve struct
	case *Resolve:
		return (*resolveCodec)(v), nil

	// Major WIT types
	case *World:
		return (*worldCodec)(v), nil
	case *Interface:
		return (*interfaceCodec)(v), nil
	case *TypeDef:
		return (*typeDefCodec)(v), nil
	case *Package:
		return (*packageCodec)(v), nil
	case *Function:
		return (*functionCodec)(v), nil

	// WIT sections
	case *[]*World:
		return (*ptrSliceCodec[World])(v), nil
	case *[]*Interface:
		return (*ptrSliceCodec[Interface])(v), nil
	case *[]*TypeDef:
		return (*ptrSliceCodec[TypeDef])(v), nil
	case *[]*Package:
		return (*ptrSliceCodec[Package])(v), nil

	// Maps
	case map[string]WorldItem:
		return mapFuncCodec[WorldItem](v), nil
	case map[string]*TypeDef:
		return mapCodec[TypeDef](v), nil

	// Callback setters of enum/interface types
	case func(WorldItem):
		return worldItemCodec(v), nil
	case func(Type):
		return &typeCodec{v, (*Resolve)(c)}, nil

	// Callback setters of concrete types
	case func(*World):
		return &indexCodec[World]{&c.Worlds, v}, nil
	case func(*Interface):
		return &indexCodec[Interface]{&c.Interfaces, v}, nil
	case func(*TypeDef):
		return &indexCodec[TypeDef]{&c.TypeDefs, v}, nil
	case func(*Package):
		return &indexCodec[Package]{&c.Packages, v}, nil
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

type resolveCodec Resolve

func (c *resolveCodec) DecodeField(name string) (any, error) {
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

type Setter[T any] func(v T)

func (s Setter[T]) Set(a any) {
	v, _ := a.(T)
	s(v)
}

type worldCodec World

func (c *worldCodec) DecodeField(name string) (any, error) {
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

type interfaceCodec Interface

func (c *interfaceCodec) DecodeField(name string) (any, error) {
	switch name {
	case "docs":
		return &c.Docs, nil
	case "name":
		return &c.Name, nil
	case "types":
		return &c.Types, nil
	case "functions":
		return &c.Functions, nil
	case "package":
		return &c.Package, nil
	}
	return nil, nil
}

type typeDefCodec TypeDef

type packageCodec Package

type typeCodec struct {
	f func(Type)
	*Resolve
}

func (c typeCodec) DecodeString(s string) error {
	var t Type
	switch s {
	// TODO
	}
	c.f(t)
	return nil
}

func (c typeCodec) DecodeInt(i int64) error {
	var t Type = realloc(&c.TypeDefs, int(i))
	c.f(t)
	return nil
}

type functionCodec Function

func (c *functionCodec) DecodeElement(name string) (any, error) {
	switch name {
	case "docs":
		return &c.Docs, nil
	case "name":
		return &c.Name, nil
	case "kind":
		return &c.Kind, nil
	case "params":
		return &c.Params, nil
	case "results":
		return &c.Results, nil
	}
	return nil, nil
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
		*s = append(*s, make([]E, i-len(*s))...)
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
