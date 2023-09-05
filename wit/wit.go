package wit

import (
	"encoding/json"
)

type Resolver interface {
	World(i int) *World
	Interface(i int) *Interface
	Type(t any) Type
	TypeDef(i int) *TypeDef
	Package(i int) *Package
}

type Resolve struct {
	Worlds     []*World     `json:"worlds"`
	Interfaces []*Interface `json:"interfaces"`
	TypeDefs   []*TypeDef   `json:"types"`
	Packages   []*Package   `json:"packages"`
}

func (r *Resolve) World(i int) *World {
	return r.Worlds[i]
}

func (r *Resolve) Interface(i int) *Interface {
	return r.Interfaces[i]
}

// TODO: should this be on TypeAny?
func (r *Resolve) Type(t any) Type {
	return nil
}

func (r *Resolve) TypeDef(i int) *TypeDef {
	return r.TypeDefs[i]
}

func (r *Resolve) Package(i int) *Package {
	return r.Packages[i]
}

func (r *Resolve) UnmarshalJSON(data []byte) error {
	// This function reads the JSON twice. The first time to get counts of
	// the worlds, interfaces, types, and packages. The second time to
	// deserialize the JSON into these structs, which contain pointers
	// to one of these values.
	var sizes struct {
		Worlds     []struct{} `json:"worlds"`
		Interfaces []struct{} `json:"interfaces"`
		TypeDefs   []struct{} `json:"types"`
		Packages   []struct{} `json:"packages"`
	}
	err := json.Unmarshal(data, &sizes)
	if err != nil {
		return err
	}

	// Allocate items
	r.Worlds = make([]*World, len(sizes.Worlds))
	for i := range r.Worlds {
		r.Worlds[i] = &World{Resolver: r}
	}
	r.Interfaces = make([]*Interface, len(sizes.Interfaces))
	for i := range r.Interfaces {
		r.Interfaces[i] = &Interface{Resolver: r}
	}
	r.TypeDefs = make([]*TypeDef, len(sizes.TypeDefs))
	for i := range r.TypeDefs {
		r.TypeDefs[i] = &TypeDef{Resolver: r}
	}
	r.Packages = make([]*Package, len(sizes.Packages))
	for i := range r.Packages {
		r.Packages[i] = &Package{Resolver: r}
	}
	return nil
}

type World struct {
	Name    string               `json:"name"`
	Docs    Docs                 `json:"docs"`
	Imports map[string]WorldItem `json:"imports"`
	Exports map[string]WorldItem `json:"exports"`
	Package *Package             `json:"package"`
	Resolver
}

func (*World) isTypeOwner() {}

type WorldItem interface {
	isWorldItem()
}

type Interface struct {
	Docs      Docs                `json:"docs"`
	Name      *string             `json:"name"`
	Types     map[string]*TypeDef `json:"types"`
	Functions map[string]Function `json:"functions"`
	Package   *Package            `json:"package"`
	Resolver
}

func (*Interface) isWorldItem() {}

func (*Interface) isTypeOwner() {}

type TypeDef struct {
	Kind  TypeDefKind `json:"kind"`
	Name  string      `json:"name,omitempty"`
	Owner TypeOwner   `json:"owner"`
	Resolver
}

func (TypeDef) isType() {}

func (*TypeDef) isWorldItem() {}

type TypeDefKind interface {
	isTypeDefKind()
}

type TypeOwner interface {
	isTypeOwner()
}

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

type Function struct {
	Docs    Docs         `json:"docs"`
	Name    string       `json:"name"`
	Kind    FunctionKind `json:"kind"`
	Params  []Param      `json:"params"`
	Results []Param      `json:"results"`
	Resolver
}

func (f *Function) UnmarshalJSON(data []byte) error {
	type function Function
	f.Kind = FunctionKindAny{f.Resolver}
	return json.Unmarshal(data, (*function)(f))
}

func (*Function) isWorldItem() {}

type FunctionKindAny struct {
	Resolver
}

func (FunctionKindAny) isFunctionKind() {}

func (u FunctionKindAny) UnmarshalJSON(data []byte) error {
	// if _, ok := u["freestanding"]; ok {
	// 	return FunctionKindFreestanding{}
	// }
	// if t, ok := u["method"]; ok {
	// 	return r.Type(t)
	// }
	return nil
}

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
	Name string `json:"contents"`
	Type Type   `json:"type"`
}

type Results interface {
	isResults()
}

type NamedResults []Param

func (NamedResults) isResults() {}

type AnonResults struct{ Type }

func (*AnonResults) isResults() {}

type Package struct {
	Name       PackageName           `json:"name"`
	Docs       Docs                  `json:"docs"`
	Interfaces map[string]*Interface `json:"interfaces"`
	Worlds     map[string]*World     `json:"worlds"`
	Resolver
}

// TODO: implement package name parsing
type PackageName string
