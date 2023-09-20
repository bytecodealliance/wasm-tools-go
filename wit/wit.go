package wit

import (
	"errors"
	"fmt"
	"strings"

	"github.com/coreos/go-semver/semver"
)

// Resolve represents a fully resolved set of WIT ([WebAssembly Interface Type])
// packages.
//
// This structure contains a graph of WIT packages and their contents
// merged together into slices organized by type. Items are sorted
// topologically and everything is fully resolved.
//
// Each item in a [Resolve] has a parent link to trace it back to the original
// package as necessary.
//
// [WebAssembly Interface Type]: https://component-model.bytecodealliance.org/wit-overview.html
type Resolve struct {
	Worlds     []*World
	Interfaces []*Interface
	TypeDefs   []*TypeDef
	Packages   []*Package
}

// A World represents all of the imports and exports of a [WebAssembly component].
//
// [WebAssembly component]: https://github.com/WebAssembly/component-model/blob/main/design/mvp/WIT.md#wit-worlds
type World struct {
	Name    string
	Imports map[string]WorldItem
	Exports map[string]WorldItem
	Package *Package
	Docs    Docs
	_typeOwner
}

// A WorldItem is any item that can be exported from or imported into a [World],
// currently either an [Interface], [TypeDef], or [Function].
type WorldItem interface{ isWorldItem() }

// _worldItem is an embeddable type that conforms to the WorldItem interface.
type _worldItem struct{}

func (_worldItem) isWorldItem() {}

// An Interface represents a [collection of types and functions], which are imported into
// or exported from a [WebAssembly component].
//
// [collection of types and functions]: https://github.com/WebAssembly/component-model/blob/main/design/mvp/WIT.md#wit-interfaces.
// [WebAssembly component]: https://github.com/WebAssembly/component-model/blob/main/design/mvp/WIT.md#wit-worlds
type Interface struct {
	Name      *string
	TypeDefs  map[string]*TypeDef
	Functions map[string]*Function
	Package   *Package
	Docs      Docs
	_worldItem
	_typeOwner
}

// TypeDef represents a WIT type definition. A TypeDef may be named or anonymous,
// and optionally belong to a [World] or [Interface].
type TypeDef struct {
	Name  *string
	Kind  TypeDefKind
	Owner TypeOwner
	Docs  Docs
	_worldItem
	_type
}

// TypeName returns the type name of t, if present.
// This partially implements the Type interface.
func (t *TypeDef) TypeName() string {
	if t.Name != nil {
		return *t.Name
	}
	return "<unnamed>"
}

// TypeDefKind represents the underlying type in a [TypeDef], which can be one of
// [Record], [Resource], [Handle], [Flags], [Tuple], [Variant], [Enum],
// [Option], [Result], [List], [Future], [Stream], or [Type].
type TypeDefKind interface{ isTypeDefKind() }

// _typeDefKind is an embeddable type that conforms to the TypeDefKind interface.
type _typeDefKind struct{}

func (_typeDefKind) isTypeDefKind() {}

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

// _primitive represents a WebAssembly Component Model [primitive type]
// mapped to its equivalent Go type.
//
// [primitive type]: https://component-model.bytecodealliance.org/wit-overview.html#primitive-types
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

// ParseType parses a WIT [primitive type] string into
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

// Package represents a [WIT package] within a [Resolve].
//
// A Package is a collection of [Interface] and [World] values. Additionally,
// a Package contains a unique identifier that affects generated components and uniquely
// identifies this particular package.
//
// [WIT package]: https://component-model.bytecodealliance.org/wit-overview.html#packages
type Package struct {
	Name       PackageName
	Interfaces map[string]*Interface
	Worlds     map[string]*World
	Docs       Docs
}

// PackageName represents a [WebAssembly Component Model] package name,
// such as [wasi:clocks@1.0.0]. It contains a namespace, name, and
// optional [SemVer] version information.
//
// [WebAssembly Component Model]: https://component-model.bytecodealliance.org/introduction.html
// [wasi:clocks@1.0.0]: https://github.com/WebAssembly/wasi-clocks
// [SemVer]: https://semver.org/
type PackageName struct {
	// Namespace specifies the package namespace, such as "wasi" in "wasi:foo/bar".
	Namespace string
	// Name specifies the kebab-name of the package.
	Name string
	// Version contains optional major/minor version information.
	Version *semver.Version
}

// ParsePackageName parses a package string into a [PackageName],
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

// Validate validates p, returning any errors.
// TODO: finish this.
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

// String implements [fmt.Stringer], returning the canonical string representation of a [PackageName].
func (pn *PackageName) String() string {
	if pn.Version == nil {
		return pn.Namespace + ":" + pn.Name
	}
	return pn.Namespace + ":" + pn.Name + "@" + pn.Version.String()
}

// Docs represent WIT documentation text extracted from comments.
type Docs struct {
	Contents string // may be empty
}
