package wit

import (
	"errors"
	"fmt"
	"strings"

	"github.com/coreos/go-semver/semver"
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
	_typeOwner
}

type WorldItem interface{ isWorldItem() }

type _worldItem struct{}

func (_worldItem) isWorldItem() {}

type Interface struct {
	Name      *string
	TypeDefs  map[string]*TypeDef
	Functions map[string]*Function
	Package   *Package
	Docs      Docs
	_worldItem
	_typeOwner
}

type TypeDef struct {
	Name  *string
	Kind  TypeDefKind
	Owner TypeOwner
	Docs  Docs
	_worldItem
	_type
}

func (t *TypeDef) TypeName() string {
	if t.Name != nil {
		return *t.Name
	}
	return "<unnamed>"
}

type TypeDefKind interface{ isTypeDefKind() }

type _typeDefKind struct{}

func (_typeDefKind) isTypeDefKind() {}

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
	_typeDefKind
}

type Field struct {
	Name string
	Type Type
	Docs Docs
}

type Resource struct{ _typeDefKind }

func (Resource) UnmarshalText() ([]byte, error) { return []byte("resource"), nil }

type Handle interface {
	isHandle()
	TypeDefKind
}

type _handle struct{ _typeDefKind }

func (_handle) isHandle() {}

type OwnHandle struct {
	Type *TypeDef
	_handle
}

type BorrowHandle struct {
	Type *TypeDef
	_handle
}

type Flags struct {
	Flags []Flag
	_typeDefKind
}

type Flag struct {
	Name string
	Docs Docs
}

type Tuple struct {
	Types []Type
	_typeDefKind
}

type Variant struct {
	Cases []Case
	_typeDefKind
}

type Case struct {
	Name string
	Type Type // Represented in Rust as Option<Type>, so Type field could be nil
	Docs Docs
}

type Enum struct {
	Cases []EnumCase
	_typeDefKind
}

type EnumCase struct {
	Name string
	Docs Docs
}

type Option struct {
	Type Type
	_typeDefKind
}

type Result struct {
	OK  Type // Represented in Rust as Option<Type>, so Type field could be nil
	Err Type // Represented in Rust as Option<Type>, so Type field could be nil
	_typeDefKind
}

type List struct {
	Type Type
	_typeDefKind
}

type Future struct {
	Type Type // Represented in Rust as Option<Type>, so Type field could be nil
	_typeDefKind
}

type Stream struct {
	Element Type // Represented in Rust as Option<Type>, so Type field could be nil
	End     Type // Represented in Rust as Option<Type>, so Type field could be nil
	_typeDefKind
}

type TypeOwner interface{ isTypeOwner() }

type _typeOwner struct{}

func (_typeOwner) isTypeOwner() {}

type Type interface {
	TypeName() string
	isType()
	TypeDefKind
}

type _type struct{ _typeDefKind }

func (_type) isType() {}

func (_type) TypeName() string { return "<unnamed>" }

// _primitive represents a WebAssembly Component Model primitive type
// mapped to its equivalent Go type.
// https://component-model.bytecodealliance.org/wit-overview.html#primitive-types
type _primitive[T any] struct{ _type }

func (_primitive[T]) isType() {}

func (_primitive[T]) TypeName() string {
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

func (t _primitive[T]) MarshalText() ([]byte, error) {
	return []byte(t.TypeName()), nil
}

type BoolType struct{ _primitive[bool] }
type S8Type struct{ _primitive[int8] }
type U8Type struct{ _primitive[uint8] }
type S16Type struct{ _primitive[int16] }
type U16Type struct{ _primitive[uint16] }
type S32Type struct{ _primitive[int32] }
type U32Type struct{ _primitive[uint32] }
type S64Type struct{ _primitive[int64] }
type U64Type struct{ _primitive[uint64] }
type Float32Type struct{ _primitive[float32] }
type Float64Type struct{ _primitive[float64] }
type CharType struct{ _primitive[char] }
type StringType struct{ _primitive[string] }

// char is defined because rune is an alias of int32
type char rune

// ParseType parses a WIT primitive type string into
// the associated Type implementation from this package.
// It returns an error if the type string is not recoginized.
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
	_worldItem
}

type FunctionKind interface {
	isFunctionKind()
}

type _functionKind struct{}

func (_functionKind) isFunctionKind() {}

type FunctionKindFreestanding struct {
	_functionKind
}

type FunctionKindMethod struct {
	Type
	_functionKind
}

type FunctionKindStatic struct {
	Type
	_functionKind
}

type FunctionKindConstructor struct {
	Type
	_functionKind
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

// PackageName represents a WebAssembly Component Model package name,
// such as `wasi:clocks@1.0.0`. It contains a namespace, name, and
// optional SemVer version information.
type PackageName struct {
	// Namespace specifies the package namespace, such as `wasi` in `wasi:foo/bar`.
	Namespace string
	// Name specifies the kebab-name of the package.
	Name string
	// Version contains optional major/minor version information.
	Version *semver.Version
}

// ParsePackageName parses a package string into a PackageName,
// returning any errors encountered. The resulting PackageName
// may not be valid.
func ParsePackageName(s string) (PackageName, error) {
	var pn PackageName
	name, ver, hasVer := strings.Cut(s, "@")
	pn.Namespace, pn.Name, _ = strings.Cut(name, ":")
	if hasVer {
		var err error
		pn.Version, err = semver.NewVersion(ver)
		if err != nil {
			return pn, err
		}
	}
	return pn, pn.Validate()
}

func (pn *PackageName) Validate() error {
	switch {
	case pn.Namespace == "":
		return errors.New("missing package namespace")
	case pn.Name == "":
		return errors.New("missing package name")
		// TODO: other validations
	}
	return nil
}

// String implements fmt.Stringer, returning the canonical string representation of a PackageName.
func (pn *PackageName) String() string {
	if pn.Version == nil {
		return pn.Namespace + ":" + pn.Name
	}
	return pn.Namespace + ":" + pn.Name + "@" + pn.Version.String()
}

type Docs struct {
	Contents string
}
