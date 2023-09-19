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
	Imports map[string]WorldItem
	Exports map[string]WorldItem
	Package *Package
	Docs    Docs
	typeOwner
}

type WorldItem interface{ isWorldItem() }

type worldItem struct{}

func (worldItem) isWorldItem() {}

type Interface struct {
	Name      *string
	TypeDefs  map[string]*TypeDef
	Functions map[string]*Function
	Package   *Package
	Docs      Docs
	worldItem
	typeOwner
}

type TypeDef struct {
	Name  *string
	Kind  TypeDefKind
	Owner TypeOwner
	Docs  Docs
	worldItem
	type_
}

func (t *TypeDef) TypeName() string {
	if t.Name != nil {
		return *t.Name
	}
	return "<unnamed>"
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
	Name string
	Type Type
	Docs Docs
}

type Resource struct{ typeDefKind }

func (Resource) UnmarshalText() ([]byte, error) { return []byte("resource"), nil }

type Handle interface {
	isHandle()
	TypeDefKind
}

type handle struct{ typeDefKind }

func (handle) isHandle() {}

type OwnHandle struct {
	Type *TypeDef
	handle
}

type BorrowHandle struct {
	Type *TypeDef
	handle
}

type Flags struct {
	Flags []Flag
	typeDefKind
}

type Flag struct {
	Name string
	Docs Docs
}

type Tuple struct {
	Types []Type
	typeDefKind
}

type Variant struct {
	Cases []Case
	typeDefKind
}

type Case struct {
	Name string
	Type Type // Represented in Rust as Option<Type>, so Type field could be nil
	Docs Docs
}

type Enum struct {
	Cases []EnumCase
	typeDefKind
}

type EnumCase struct {
	Name string
	Docs Docs
}

type Option struct {
	Type Type
	typeDefKind
}

type Result struct {
	OK  Type // Represented in Rust as Option<Type>, so Type field could be nil
	Err Type // Represented in Rust as Option<Type>, so Type field could be nil
	typeDefKind
}

type List struct {
	Type Type
	typeDefKind
}

type Future struct {
	Type Type // Represented in Rust as Option<Type>, so Type field could be nil
	typeDefKind
}

type Stream struct {
	Element Type // Represented in Rust as Option<Type>, so Type field could be nil
	End     Type // Represented in Rust as Option<Type>, so Type field could be nil
	typeDefKind
}

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

func (type_) TypeName() string { return "<unnamed>" }

// primitiveType represents a WebAssembly Component Model primitive type
// mapped to its equivalent Go type.
// https://component-model.bytecodealliance.org/wit-overview.html#primitive-types
type primitiveType[T any] struct{ type_ }

func (primitiveType[T]) isType() {}

func (primitiveType[T]) TypeName() string {
	var v T
	switch any(v).(type) {
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
	return "<undefined>"
}

func (t primitiveType[T]) MarshalText() ([]byte, error) {
	return []byte(t.TypeName()), nil
}

type BoolType struct{ primitiveType[bool] }
type S8Type struct{ primitiveType[int8] }
type U8Type struct{ primitiveType[uint8] }
type S16Type struct{ primitiveType[int16] }
type U16Type struct{ primitiveType[uint16] }
type S32Type struct{ primitiveType[int32] }
type U32Type struct{ primitiveType[uint32] }
type S64Type struct{ primitiveType[int64] }
type U64Type struct{ primitiveType[uint64] }
type Float32Type struct{ primitiveType[float32] }
type Float64Type struct{ primitiveType[float64] }
type CharType struct{ primitiveType[char] }
type StringType struct{ primitiveType[string] }

// char is defined because rune is an alias of int32
type char int32

func ParseType(s string) (Type, error) {
	switch s {
	case "bool":
		return BoolType{}, nil
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
	Name    string
	Kind    FunctionKind
	Params  []Param
	Results []Param
	Docs    Docs
	worldItem
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
	Interfaces map[string]*Interface
	Worlds     map[string]*World
	Docs       Docs
}

// TODO: implement package name parsing
type PackageName string

type Docs struct {
	Contents string
}
