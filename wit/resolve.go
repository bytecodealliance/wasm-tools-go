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
// This partially implements the [Type] interface.
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

// Record represents a WIT [record type], akin to a struct.
//
// [record type]: https://component-model.bytecodealliance.org/wit-overview.html#records
type Record struct {
	Fields []Field
	_typeDefKind
}

// Field represents a field in a [Record].
type Field struct {
	Name string
	Type Type
	Docs Docs
}

// Resource represents a WIT resource type.
type Resource struct{ _typeDefKind }

func (Resource) UnmarshalText() ([]byte, error) { return []byte("resource"), nil }

// Handle represents a WIT handle type.
type Handle interface {
	isHandle()
	TypeDefKind
}

type _handle struct{ _typeDefKind }

func (_handle) isHandle() {}

// OwnHandle represents an owned WIT handle type.
type OwnHandle struct {
	Type *TypeDef
	_handle
}

// BorrowHandle represents an borrowed WIT handle type.
type BorrowHandle struct {
	Type *TypeDef
	_handle
}

// Flags represents a WIT [flags type], stored as a bitfield.
//
// [flags type]: https://component-model.bytecodealliance.org/wit-overview.html#flags
type Flags struct {
	Flags []Flag
	_typeDefKind
}

// Flag represents a single flag value in a [Flags] type.
type Flag struct {
	Name string
	Docs Docs
}

// Tuple represents a WIT [tuple type].
// A tuple type is an ordered fixed length sequence of values of specified types.
// It is similar to a [Record], except that the fields are identified by their order instead of by names.
//
// [tuple type]: https://component-model.bytecodealliance.org/wit-overview.html#tuples
type Tuple struct {
	Types []Type
	_typeDefKind
}

// Variant represents a WIT [variant type], a tagged/discriminated union.
// A variant type declares one or more cases. Each case has a name and, optionally,
// a type of data associated with that case.
//
// [variant type]: https://component-model.bytecodealliance.org/wit-overview.html#variants
type Variant struct {
	Cases []Case
	_typeDefKind
}

// Case represents a single case in a [Variant].
type Case struct {
	Name string
	Type Type // Represented in Rust as Option<Type>, so Type field could be nil
	Docs Docs
}

// Enum represents a WIT [enum type], which is a [Variant] without associated data.
// The equivalent in Go is a set of const identifiers declared with iota.
//
// [enum type]: https://component-model.bytecodealliance.org/wit-overview.html#enums
type Enum struct {
	Cases []EnumCase
	_typeDefKind
}

// EnumCase represents a single case in an [Enum].
type EnumCase struct {
	Name string
	Docs Docs
}

// Option represents a WIT [option type], a special case of [Variant]. An Option can
// contain a value of a single type, either build-in or user defined, or no value.
// The equivalent in Go for an option<string> could be represented as *string.
//
// [option type]: https://component-model.bytecodealliance.org/wit-overview.html#options
type Option struct {
	Type Type
	_typeDefKind
}

// Result represents a WIT [result type], which is the result of a function call,
// returning an optional value and/or an optional error. It is roughly equivalent to
// the Go pattern of returning (T, error).
//
// [result type]: https://component-model.bytecodealliance.org/wit-overview.html#results
type Result struct {
	OK  Type // Represented in Rust as Option<Type>, so Type field could be nil
	Err Type // Represented in Rust as Option<Type>, so Type field could be nil
	_typeDefKind
}

// List represents a WIT [list type], which is an ordered vector of an arbitrary type.
//
// [list type]: https://component-model.bytecodealliance.org/wit-overview.html#lists
type List struct {
	Type Type
	_typeDefKind
}

// Future represents a WIT [future type], expected to be part of [WASI Preview 3].
//
// [future type]: https://github.com/bytecodealliance/wit-bindgen/issues/270
// [WASI Preview 3]: https://bytecodealliance.org/articles/webassembly-the-updated-roadmap-for-developers
type Future struct {
	Type Type // Represented in Rust as Option<Type>, so Type field could be nil
	_typeDefKind
}

// Stream represents a WIT [stream type], expected to be part of [WASI Preview 3].
//
// [stream type]: https://github.com/WebAssembly/WASI/blob/main/docs/WitInWasi.md#streams
// [WASI Preview 3]: https://bytecodealliance.org/articles/webassembly-the-updated-roadmap-for-developers
type Stream struct {
	Element Type // Represented in Rust as Option<Type>, so Type field could be nil
	End     Type // Represented in Rust as Option<Type>, so Type field could be nil
	_typeDefKind
}

// TypeOwner is the interface implemented by any type that can own a TypeDef,
// currently [World] and [Interface].
type TypeOwner interface{ isTypeOwner() }

type _typeOwner struct{}

func (_typeOwner) isTypeOwner() {}

// Type is the interface implemented by any type definition. This can be a
// [primitive type] or a user-defined type in a [TypeDef].
//
// [primitive type]: https://component-model.bytecodealliance.org/wit-overview.html#primitive-types
type Type interface {
	TypeName() string
	isType()
	TypeDefKind
}

// _type is an embeddable type that conforms to the Type interface.
type _type struct{ _typeDefKind }

func (_type) isType() {}

func (_type) TypeName() string { return "<unnamed>" }

// Primitive is a type constriant of the Go equivalents of WIT [primitive types].
//
// [primitive types]: https://component-model.bytecodealliance.org/wit-overview.html#primitive-types
type Primitive interface {
	bool | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64 | float32 | float64 | char | string
}

// char is defined because [rune] is an alias of [int32]
type char rune

// _primitive represents a WebAssembly Component Model [primitive type]
// mapped to its equivalent Go type.
//
// [primitive type]: https://component-model.bytecodealliance.org/wit-overview.html#primitive-types
type _primitive[T Primitive] struct{ _type }

// _primitive is a generic embeddable type that conforms to the [Type] interface.
func (_primitive[T]) isType() {}

// TypeName partially implements the [Type] interface.
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

// Function represents a WIT [function].
// Functions can be freestanding, methods, constructors or static.
//
// [function]: https://component-model.bytecodealliance.org/wit-overview.html#functions
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
