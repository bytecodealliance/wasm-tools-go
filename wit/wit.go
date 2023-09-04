package wit

type Resolve struct {
	Worlds     []*World
	Interfaces []*Interface
	Types      []*TypeDef
	Packages   []*Package
}

type World struct {
	Name    string               `json:"name"`
	Docs    Docs                 `json:"docs"`
	Imports map[string]WorldItem `json:"imports"`
	Exports map[string]WorldItem `json:"exports"`
	Package *Package             `json:"package"`
}

func (*World) isTypeOwner() {}

type WorldItem interface {
	isWorldItem()
}

type Package struct {
	Name       PackageName           `json:"name"`
	Docs       Docs                  `json:"docs"`
	Interfaces map[string]*Interface `json:"interfaces"`
	Worlds     map[string]*World     `json:"worlds"`
}

// TODO: implement package name parsing
type PackageName string

type Interface struct {
	Docs      Docs                `json:"docs"`
	Name      *string             `json:"name"`
	Types     map[string]*TypeDef `json:"types"`
	Functions map[string]Function `json:"functions"`
	Package   *Package            `json:"package"`
}

func (*Interface) isWorldItem() {}

func (*Interface) isTypeOwner() {}

type TypeDef struct {
	Kind  TypeDefKind `json:"kind"`
	Name  string      `json:"name,omitempty"`
	Owner TypeOwner   `json:"owner"`
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
	Docs    Docs
	Name    string
	Kind    FunctionKind
	Params  []Param
	Results []Param
}

func (*Function) isWorldItem() {}

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
	Type Type
}

type Results interface {
	isResults()
}

type NamedResults []Param

func (NamedResults) isResults() {}

type AnonResults struct{ Type }

func (*AnonResults) isResults() {}
