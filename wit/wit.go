package wit

import (
	"fmt"
	"strconv"
	"unsafe"

	"github.com/ydnar/wit-bindgen-go/internal/codec"
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
	isType()
	TypeDefKind
}

type type_ struct{ typeDefKind }

func (type_) isType() {}

type IntType interface {
	isIntType()
	Type
}

type intType struct{ type_ }

func (intType) isIntType() {}

/*
func (intType[T]) MarshalText() ([]byte, error) {
	var v T
	v -= 1
	pfx := "u"
	if v < 0 {
		pfx = "s"
	}
	return []byte(pfx + string(unsafe.Sizeof(v))), nil
}
*/

type sintType[T codec.Signed] struct{ intType }

func (sintType[T]) MarshalText() ([]byte, error) {
	var v T
	return []byte("s" + strconv.Itoa(int(unsafe.Sizeof(v)*8))), nil
}

type uintType[T codec.Unsigned] struct{ intType }

func (uintType[T]) MarshalText() ([]byte, error) {
	var v T
	return []byte("u" + strconv.Itoa(int(unsafe.Sizeof(v)*8))), nil
}

type FloatType interface {
	isFloatType()
	Type
}

type floatType[T codec.Float] struct{ type_ }

func (floatType[T]) isFloatType() {}

func (floatType[T]) MarshalText() ([]byte, error) {
	var v T
	return []byte("float" + strconv.Itoa(int(unsafe.Sizeof(v)*8))), nil
}

type BoolType struct{ type_ }

func (BoolType) MarshalText() ([]byte, error) { return []byte("bool"), nil }

type S8Type struct{ sintType[int8] }
type U8Type struct{ uintType[uint8] }
type S16Type struct{ sintType[int16] }
type U16Type struct{ uintType[uint16] }
type S32Type struct{ sintType[int32] }
type U32Type struct{ uintType[uint32] }
type S64Type struct{ sintType[int64] }
type U64Type struct{ uintType[uint64] }
type Float32Type struct{ floatType[float32] }
type Float64Type struct{ floatType[float64] }
type CharType struct{ type_ }

func (CharType) MarshalText() ([]byte, error) { return []byte("char"), nil }

type StringType struct{ type_ }

func (StringType) MarshalText() ([]byte, error) { return []byte("string"), nil }

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
