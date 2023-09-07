package wit

type Resolve struct {
	Worlds     []*World
	Interfaces []*Interface
	TypeDefs   []*TypeDef
	Packages   []*Package
}

type World struct {
	Name    string
	Docs    Docs
	Imports map[string]WorldItem
	Exports map[string]WorldItem
	Package *Package
}

func (*World) isTypeOwner() {}

type WorldItem interface {
	isWorldItem()
}

type Interface struct {
	Docs      Docs
	Name      *string
	Types     map[string]*TypeDef
	Functions map[string]*Function
	Package   *Package
}

func (*Interface) isWorldItem() {}

func (*Interface) isTypeOwner() {}

type TypeDef struct {
	Kind  TypeDefKind
	Name  string
	Owner TypeOwner
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
	Contents *string
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

type Package struct {
	Name       PackageName
	Docs       Docs
	Interfaces map[string]*Interface
	Worlds     map[string]*World
}

// TODO: implement package name parsing
type PackageName string
