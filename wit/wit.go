package wit

type Resolve struct {
	Worlds     []*World
	Interfaces []*Interface
	Types      []*TypeDef
	Packages   []*Package
}

// worldAt returns a pointer to the worldAt at index i.
// The underlying slice will be lengthened if i is out of bounds.
// Returns nil if i < 0.
func (r *Resolve) worldAt(i int) *World {
	return element(&r.Worlds, i)
}

// Interface returns a pointer to the Interface at index i.
// The underlying slice will be lengthened if i is out of bounds.
// Returns nil if i < 0.
func (r *Resolve) Interface(i int) *Interface {
	return element(&r.Interfaces, i)
}

// Type returns a pointer to the TypeDef at index i.
// The underlying slice will be lengthened if i is out of bounds.
// Returns nil if i < 0.
func (r *Resolve) Type(i int) *TypeDef {
	return element(&r.Types, i)
}

// Package returns a pointer to the Package at index i.
// The underlying slice will be lengthened if i is out of bounds.
// Returns nil if i < 0.
func (r *Resolve) Package(i int) *Package {
	return element(&r.Packages, i)
}

// element returns the value of slice s at index i,
// reallocating the slice if necessary. s must be a slice
// of pointers, because the underlying backing to s might
// change when reallocated.
// If the value at s[i] is nil, a new *E will be allocated.
func element[S ~*[]*E, E any](s S, i int) *E {
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

type World struct {
	Name string `json:"name"`
	Docs Docs   `json:"docs"`
}

type Package struct {
	Name       PackageName `json:"name"`
	Docs       Docs        `json:"docs"`
	Interfaces map[string]*Interface
}

// TODO: implement package name parsing
type PackageName string

type WorldItem interface {
	isWorldItem()
}

type Interface struct {
	Docs      Docs                `json:"docs"`
	Name      *string             `json:"name"`
	Types     map[string]*TypeDef `json:"types"`
	Functions map[string]Function `json:"functions"`
	Package   *Package            `json:"package"`
}

func (i *Interface) isWorldItem() {}

/*
pub enum Type {
    Bool,
    U8,
    U16,
    U32,
    U64,
    S8,
    S16,
    S32,
    S64,
    Float32,
    Float64,
    Char,
    String,
    Id(TypeId),
}
*/

type Type interface{ isType() }

type BoolType struct{}

func (BoolType) isType() {}

type U8Type struct{}

func (U8Type) isType() {}

// TODO: rest of the types

type TypeDef struct {
	// TODO
}

func (TypeDef) isType() {}

type Function struct {
	Docs    Docs
	Name    string
	Kind    FunctionKind
	Params  []Param // Vec<(String, Type)>;
	Results Results // enum
}

func (Function) isWorldItem() {}

type FunctionKind interface{ isFunctionKind() }

type FunctionKindFreestanding struct{}

func (FunctionKindFreestanding) isFunctionKind() {}

type FunctionKindMethod struct{ Type }

func (FunctionKindMethod) isFunctionKind() {}

type FunctionKindStatic struct{ Type }

func (FunctionKindStatic) isFunctionKind() {}

type FunctionKindConstructor struct{ Type }

func (FunctionKindConstructor) isFunctionKind() {}

type Docs struct {
	Contents *string `json:"contents"`
}

type Param struct {
	Name string
	Type *Type
}

type Results struct {
	Named []Param
	Anon  *Type
}

func init() {
}
