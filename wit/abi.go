package wit

// Align aligns ptr with alignment align.
func Align(ptr, align uintptr) uintptr {
	return (ptr + align - 1) &^ (align - 1)
}

// Discriminant returns the smallest WIT integer type that can represent 0...n.
// Used by the [Canonical ABI] for [Variant] types.
//
// [Canonical ABI]: https://github.com/WebAssembly/component-model/blob/main/design/mvp/CanonicalABI.md#alignment
func Discriminant(n int) Type {
	switch {
	case n <= 1<<8:
		return U8{}
	case n <= 1<<16:
		return U16{}
	}
	return U32{}
}

// ABI is the interface implemented by any type that reports its [Canonical ABI] [size], [alignment],
// whether the type contains a pointer (e.g. [List] or [String]), and its [flat] ABI representation.
//
// [Canonical ABI]: https://github.com/WebAssembly/component-model/blob/main/design/mvp/CanonicalABI.md
// [size]: https://github.com/WebAssembly/component-model/blob/main/design/mvp/CanonicalABI.md#size
// [alignment]: https://github.com/WebAssembly/component-model/blob/main/design/mvp/CanonicalABI.md#alignment
// [flat]: https://github.com/WebAssembly/component-model/blob/main/design/mvp/CanonicalABI.md#flattening
type ABI interface {
	Size() uintptr
	Align() uintptr
	HasPointer() bool
	Flat() []Type
}

// Despecializer is the interface implemented by any [TypeDefKind] that can
// [despecialize] itself into another TypeDefKind. Examples include [Result],
// which despecializes into a [Variant] with two cases, "ok" and "error".
// See the [canonical ABI documentation] for more information.
//
// [despecialize]: https://github.com/WebAssembly/component-model/blob/main/design/mvp/CanonicalABI.md#despecialization
// [canonical ABI documentation]: https://github.com/WebAssembly/component-model/blob/main/design/mvp/CanonicalABI.md#despecialization
type Despecializer interface {
	Despecialize() TypeDefKind
}

// Despecialize [despecializes] k if k implements [Despecializer].
// Otherwise, it returns k unmodified.
// See the [canonical ABI documentation] for more information.
//
// [despecializes]: https://github.com/WebAssembly/component-model/blob/main/design/mvp/CanonicalABI.md#despecialization
// [canonical ABI documentation]: https://github.com/WebAssembly/component-model/blob/main/design/mvp/CanonicalABI.md#despecialization
func Despecialize(k TypeDefKind) TypeDefKind {
	if d, ok := k.(Despecializer); ok {
		return d.Despecialize()
	}
	return k
}

// Op represents the [Canonical ABI] [lift] and [lower] operations, for lowering into or lifting out of linear memory.
//
// [Canonical ABI]: https://github.com/WebAssembly/component-model/blob/main/design/mvp/Explainer.md#canonical-abi
// [lift]: https://github.com/WebAssembly/component-model/blob/main/design/mvp/CanonicalABI.md#canon-lift
// [lower]: https://github.com/WebAssembly/component-model/blob/main/design/mvp/CanonicalABI.md#canon-lower
type Op uint8

const (
	// Lift represents the Canonical ABI [lift] direction, lifting Component Model types out of linear memory.
	// Used for exporting functions to the WebAssembly host (or another component) using //go:wasmexport.
	//
	// [lift]: https://github.com/WebAssembly/component-model/blob/main/design/mvp/CanonicalABI.md#canon-lift
	Lift Op = 0

	// Lower represents the Canonical ABI [lower] operation, lowering Component Model types into linear memory.
	// Used for calling functions imported from the WebAssembly host (or another component) using //go:wasmimport.
	//
	// [lower]: https://github.com/WebAssembly/component-model/blob/main/design/mvp/CanonicalABI.md#canon-lower
	Lower Op = 1

	// MaxFlatParams is the maximum number of [flattened parameters] a function can have
	// as defined in the Component Model Canonical ABI.
	//
	// [flattened parameters]: https://github.com/WebAssembly/component-model/blob/main/design/mvp/CanonicalABI.md#flattening
	MaxFlatParams = 16

	// MaxFlatResults is the maximum number of [flattened results] a function can have
	// as defined in the Component Model Canonical ABI.
	//
	// [flattened results]: https://github.com/WebAssembly/component-model/blob/main/design/mvp/CanonicalABI.md#flattening
	MaxFlatResults = 1
)

// ResourceDrop returns the implied [resource-drop] method for t.
// If t is not a [Resource], this returns nil.
//
// [resource-drop]: https://github.com/WebAssembly/component-model/blob/main/design/mvp/CanonicalABI.md#canon-resourcedrop
func (t *TypeDef) ResourceDrop() *Function {
	if _, ok := t.Kind.(*Resource); !ok {
		return nil
	}
	f := &Function{
		Name:   "[resource-drop]" + t.TypeName(),
		Kind:   &Method{Type: t},
		Params: []Param{{Name: "self", Type: t}},
		Docs:   Docs{Contents: "Drops a resource handle."},
	}
	return f
}

// CoreFunction returns a [Core WebAssembly function] of [Function] f.
// Its params and results may be [flattened] according to the Canonical ABI specification.
// The flattening rules vary based on whether the returned function is imported or exported,
// e.g. using go:wasmimport or go:wasmexport.
//
// [Core WebAssembly function]: https://webassembly.github.io/spec/core/syntax/modules.html#syntax-func
// [flattened]: https://github.com/WebAssembly/component-model/blob/main/design/mvp/CanonicalABI.md#flattening
func (f *Function) CoreFunction(op Op) *Function {
	if len(f.Params) == 0 && len(f.Results) == 0 {
		return f
	}

	// Clone the function
	cf := *f

	// Max 16 params
	if len(flatParams(f.Params)) > MaxFlatParams {
		cf.Params = []Param{compoundParam("param", "params", f.Params)}
	}

	// Max 1 result
	if len(flatParams(f.Results)) > MaxFlatResults {
		p := compoundParam("result", "results", f.Results)
		if op == Lift {
			cf.Results = []Param{p}
		} else {
			cf.Params = append(cf.Params, p)
			cf.Results = nil
		}
	}

	return &cf
}

func flatParams(params []Param) []Type {
	flat := make([]Type, 0, len(params))
	for _, p := range params {
		flat = append(flat, p.Type.Flat()...)
	}
	return flat
}

// compoundParam returns a single param that represents
// the combined param(s), using a [Pointer].
func compoundParam(singular, plural string, params []Param) Param {
	if len(params) == 0 {
		panic("BUG: len(params) == 0")
	}

	name := params[0].Name
	var t Type

	if len(params) == 1 {
		if name == "" {
			name = singular
		}
		t = params[0].Type
	} else {
		name = plural
		r := &Record{}
		t = &TypeDef{Kind: r}
		for _, p := range params {
			r.Fields = append(r.Fields,
				Field{
					Name: p.Name,
					Type: p.Type,
				})
		}
	}

	return Param{
		Name: name,
		Type: &TypeDef{Kind: &Pointer{Type: t}},
	}
}
