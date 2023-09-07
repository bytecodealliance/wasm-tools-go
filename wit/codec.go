package wit

type Codec Resolve

func (c *Codec) Codec(v any) (any, error) {
	switch v := v.(type) {
	case *Resolve:
		return (*resolveCodec)(v), nil

	case *World:
		return (*worldCodec)(v), nil
	case *Interface:
		return (*interfaceCodec)(v), nil
	case *TypeDef:
		return (*typeDefCodec)(v), nil
	case *Package:
		return (*packageCodec)(v), nil

	case *[]*World:
		return (*sliceCodec[World])(v), nil
	case *[]*Interface:
		return (*sliceCodec[Interface])(v), nil
	case *[]*TypeDef:
		return (*sliceCodec[TypeDef])(v), nil
	case *[]*Package:
		return (*sliceCodec[Package])(v), nil

	case map[string]WorldItem:
		return mapFuncCodec[WorldItem](v), nil
	case map[string]*TypeDef:
		return mapCodec[TypeDef](v), nil
	}
	return nil, nil
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

func (c typeCodec) DecodeValue(v string) error {
	var t Type
	c.f(t)
	return nil
}

type sliceCodec[T any] []*T

func (c *sliceCodec[T]) DecodeElement(i int) (any, error) {
	return remake(c, i), nil
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
