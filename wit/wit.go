package wit

import "fmt"

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

type type_ struct{}

func (type_) isType() {}

type IntType interface {
	isIntType()
	Type
}

type intType struct{ type_ }

func (intType) isIntType() {}

type FloatType interface {
	isFloatType()
	Type
}

type floatType struct{ type_ }

func (floatType) isFloatType() {}

type BoolType struct{ type_ }
type S8Type struct{ intType }
type U8Type struct{ intType }
type S16Type struct{ intType }
type U16Type struct{ intType }
type S32Type struct{ intType }
type U32Type struct{ intType }
type S64Type struct{ intType }
type U64Type struct{ intType }
type Float32Type struct{ floatType }
type Float64Type struct{ floatType }
type CharType struct{ type_ }
type StringType struct{ type_ }

func ParseType(s string) (Type, error) {
	switch s {
	case "s8":
		return S8Type{}, nil
	case "u8":
		return U8Type{}, nil
	case "s16":
		return S16Type{}, nil
	case "u16":
		return U16Type{}, nil
	case "s32":
		return S32Type{}, nil
	case "u32":
		return U32Type{}, nil
	case "s64":
		return S64Type{}, nil
	case "u64":
		return U64Type{}, nil
	case "float32":
		return Float32Type{}, nil
	case "float64":
		return Float64Type{}, nil
	case "char":
		return CharType{}, nil
	case "string":
		return StringType{}, nil
	}
	return nil, fmt.Errorf("unknown type %q", s)
}

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

type Package struct {
	Name       PackageName
	Docs       Docs
	Interfaces map[string]*Interface
	Worlds     map[string]*World
}

// TODO: implement package name parsing
type PackageName string
