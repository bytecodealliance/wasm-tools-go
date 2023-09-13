package wit

import (
	"fmt"
)

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
	typeOwner
}

type WorldItem interface{ isWorldItem() }

type worldItem struct{}

func (worldItem) isWorldItem() {}

type Interface struct {
	Docs      Docs
	Name      *string
	TypeDefs  map[string]*TypeDef
	Functions map[string]*Function
	Package   *Package `json:"-"`
	worldItem
	typeOwner
}

type TypeDef struct {
	Kind  TypeDefKind
	Name  string
	Owner TypeOwner `json:"-"`
	worldItem
	type_
}

func (t *TypeDef) TypeName() string {
	return t.Name
}

type TypeDefKind interface{ isTypeDefKind() }

type typeDefKind struct{}

func (typeDefKind) isTypeDefKind() {}

/*
	Record(Record),
	Resource,
	Handle(Handle),
	Flags(Flags),
	Tuple(Tuple),
	Variant(Variant),
	Enum(Enum),
	Option(Type),
	Result(Result_),
	List(Type),
	Future(Option<Type>),
	Stream(Stream),
	Type(Type),
*/

type Record struct {
	Fields []Field
	typeDefKind
}

type Field struct {
	Docs Docs
	Name string
	Type Type
}

type Resource struct {
	// TODO
	typeDefKind
}

type Handle struct {
	// TODO
	typeDefKind
}

type Flags struct {
	// TODO
	typeDefKind
}

type Tuple struct {
	// TODO
	typeDefKind
}

type Variant struct {
	// TODO
	typeDefKind
}

type Enum struct {
	// TODO
	typeDefKind
}

type Option struct{ Type }

type Result struct {
	// TODO
	typeDefKind
}

type List struct {
	// TODO
	typeDefKind
}

type Future struct{ Type } // Represented in Rust as Option<Type>, so Type field could be nil

type Stream struct{ Type }

type TypeOwner interface{ isTypeOwner() }

type typeOwner struct{}

func (typeOwner) isTypeOwner() {}

type Type interface {
	TypeName() string
	isType()
	TypeDefKind
}

type type_ struct{ typeDefKind }

func (type_) isType() {}

func (type_) TypeName() string { return "<undefined>" }

type coreType[T any] struct{ type_ }

func (coreType[T]) isType() {}

func (coreType[T]) TypeName() string {
	var v T
	switch any(v).(type) {
	case nil:
		return "nil"
	case bool:
		return "bool"
	case int8:
		return "s8"
	case uint8:
		return "u8"
	case int16:
		return "s16"
	case uint16:
		return "u16"
	case int32:
		return "s32"
	case uint32:
		return "u32"
	case int64:
		return "s64"
	case uint64:
		return "u64"
	case float32:
		return "float32"
	case float64:
		return "float64"
	case char:
		return "char"
	case string:
		return "string"
	}
	return ""
}

func (t coreType[T]) MarshalText() ([]byte, error) {
	return []byte(t.TypeName()), nil
}

type BoolType struct{ coreType[bool] }
type S8Type struct{ coreType[int8] }
type U8Type struct{ coreType[uint8] }
type S16Type struct{ coreType[int16] }
type U16Type struct{ coreType[uint16] }
type S32Type struct{ coreType[int32] }
type U32Type struct{ coreType[uint32] }
type S64Type struct{ coreType[int64] }
type U64Type struct{ coreType[uint64] }
type Float32Type struct{ coreType[float32] }
type Float64Type struct{ coreType[float64] }
type CharType struct{ coreType[char] }
type StringType struct{ coreType[string] }

// char is defined because rune is an alias of int32
type char int32

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

type Function struct {
	Docs    Docs
	Name    string
	Kind    FunctionKind
	Params  []Param
	Results []Param
}

type FunctionKind interface {
	isFunctionKind()
}

type functionKind struct{}

func (functionKind) isFunctionKind() {}

type FunctionKindFreestanding struct {
	functionKind
}

type FunctionKindMethod struct {
	Type
	functionKind
}

type FunctionKindStatic struct {
	Type
	functionKind
}

type FunctionKindConstructor struct {
	Type
	functionKind
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

type Docs struct {
	Contents *string
}
