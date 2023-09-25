package wit

import (
	"fmt"
	"strings"

	"github.com/ydnar/wasm-tools-go/internal/codec"
)

// Node is the interface implemented by the WIT ([WebAssembly Interface Type])
// types in this package.
//
// [WebAssembly Interface Type]: https://component-model.bytecodealliance.org/wit-overview.html
type Node interface {
	WIT(ctx Node, name string) string
}

type _node struct{}

func (_node) WIT(ctx Node, name string) string { return fmt.Sprintf("// TODO(%s)", name) }

func indent(s string) string {
	const ws = "    "
	return strings.TrimSuffix(ws+strings.Replace(s, "\n", "\n"+ws, -1), ws)
}

// WIT returns the WIT representation of r.
func (r *Resolve) WIT(_ Node, _ string) string {
	b := &strings.Builder{}
	for i, p := range r.Packages {
		if i > 0 {
			b.WriteRune('\n')
			b.WriteRune('\n')
		}
		b.WriteString(p.WIT(r, ""))
	}
	return b.String()
}

// WIT returns the WIT representation of w.
func (w *World) WIT(ctx Node, name string) string {
	var b strings.Builder
	// TODO: docs
	fmt.Fprintf(&b, "world %s {", name) // TODO: compare to w.Name?
	if len(w.Imports) > 0 || len(w.Exports) > 0 {
		b.WriteRune('\n')
		for _, name := range codec.SortedKeys(w.Imports) {
			b.WriteString(indent(w.itemWIT("import", name, w.Imports[name])))
			b.WriteRune('\n')
		}
		for _, name := range codec.SortedKeys(w.Exports) {
			b.WriteString(indent(w.itemWIT("export", name, w.Exports[name])))
			b.WriteRune('\n')
		}
	}
	b.WriteRune('}')
	return b.String()
}

func (w *World) itemWIT(motion, name string, v WorldItem) string {
	switch v := v.(type) {
	case *Interface, *Function:
		return motion + " " + v.WIT(w, name)
	case *TypeDef:
		return v.WIT(w, name) // no motion, in Imports only
	}
	panic("BUG: unknown WorldItem")
}

// WIT returns the WIT representation of i.
func (i *Interface) WIT(ctx Node, name string) string {
	if i.Name != nil && name == "" {
		name = *i.Name
	}

	var b strings.Builder

	// TODO: docs

	switch ctx := ctx.(type) {
	case *Package:
		b.WriteString("interface ")
		b.WriteString(name)
		b.WriteRune(' ')
	case *World:
		if i.Package != ctx.Package {
			// Import by name from another package
			// TODO: check i.Name != nil
			return fmt.Sprintf("%s/%s", i.Package.Name.String(), *i.Name)
		} else if i.Name != nil {
			// Import by name within same package
			return *i.Name
		}
		// Otherwise, this is an inline interface decl
		b.WriteString("interface ")
		b.WriteString(name)
		b.WriteRune(' ')
	}

	b.WriteRune('{')
	if len(i.TypeDefs) > 0 || len(i.Functions) > 0 {
		b.WriteRune('\n')
		n := 0
		for _, name := range codec.SortedKeys(i.TypeDefs) {
			// if n > 0 {
			// 	b.WriteRune('\n')
			// }
			b.WriteString(indent(i.TypeDefs[name].WIT(i, name)))
			b.WriteRune('\n')
			n++
		}
		for _, name := range codec.SortedKeys(i.Functions) {
			// if n > 0 {
			// 	b.WriteRune('\n')
			// }
			b.WriteString(indent(i.Functions[name].WIT(i, name)))
			b.WriteRune('\n')
			n++
		}
	}
	b.WriteRune('}')
	return b.String()
}

// WIT returns the WIT representation of [TypeDef] t.
func (t *TypeDef) WIT(ctx Node, name string) string {
	if t.Name != nil && name == "" {
		name = *t.Name
	}
	switch ctx := ctx.(type) {
	// If context is another TypeDef, then this is an imported type.
	case *TypeDef:
		// Emit an type alias if same Owner.
		if t.Owner == ctx.Owner && t.Name != nil {
			return "type " + name + " = " + *t.Name
		}

		// TODO: add a TypeOwnerName method to TypeDef.
		var ownerName string
		var pkg *Package
		switch owner := t.Owner.(type) {
		case *Interface:
			ownerName = *owner.Name
			pkg = owner.Package
		case *World:
			ownerName = owner.Name
			pkg = owner.Package
		}
		// TODO: use less-qualified name (without package) if this is an import within the same package.
		if t.Name != nil && *t.Name != name {
			return fmt.Sprintf("use %s/%s.{%s as %s}", pkg.Name.String(), ownerName, *t.Name, name)
		}
		return fmt.Sprintf("use %s/%s.{%s}", pkg.Name.String(), ownerName, name)

	case *World, *Interface:
		switch t.Kind.(type) {
		case *TypeDef:
			return t.Kind.WIT(t, name)
		}
		return t.Kind.WIT(ctx, name)
	}
	if name != "" {
		return name
	}
	return t.Kind.WIT(ctx, name)
}

func (r *Record) WIT(ctx Node, name string) string {
	var b strings.Builder
	b.WriteString("record ")
	b.WriteString(name)
	b.WriteString(" {")
	if len(r.Fields) > 0 {
		b.WriteRune('\n')
		for i := range r.Fields {
			if i > 0 {
				b.WriteString(",\n")
			}
			b.WriteString(indent(r.Fields[i].WIT(ctx, "")))
		}
		b.WriteRune('\n')
	}
	b.WriteString("}")
	return b.String()
}

func (f *Field) WIT(ctx Node, name string) string {
	return f.Name + ": " + f.Type.WIT(f, "")
}

func (t *Tuple) WIT(ctx Node, _ string) string {
	var b strings.Builder
	b.WriteString("tuple<")
	for i := range t.Types {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(t.Types[i].WIT(t, ""))
	}
	b.WriteString(">")
	return b.String()
}

func (l *List) WIT(ctx Node, _ string) string {
	return "list<" + l.Type.WIT(l, "") + ">"
}

// WIT returns the WIT representation of [primitive type] T.
//
// [primitive type]: https://component-model.bytecodealliance.org/wit-overview.html#primitive-types
func (p _primitive[T]) WIT(ctx Node, name string) string {
	if name != "" {
		return "type " + name + " = " + p.String()
	}
	return p.String()
}

// WIT returns the WIT representation of f.
func (f *Function) WIT(ctx Node, name string) string {
	var b strings.Builder
	// TODO: docs
	b.WriteString(name)
	b.WriteString(": func(")
	b.WriteString(paramsWIT(f.Params))
	b.WriteRune(')')
	if len(f.Results) > 0 {
		b.WriteString(" -> ")
		b.WriteString(paramsWIT(f.Results))
	}
	return b.String()
}

func paramsWIT(params []Param) string {
	var b strings.Builder
	for i, param := range params {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(param.WIT(nil, ""))
	}
	return b.String()
}

func (p *Param) WIT(_ Node, _ string) string {
	if p.Name == "" {
		return p.Type.WIT(p, "")
	}
	return p.Name + ": " + p.Type.WIT(p, "")
}

func (p *Package) WIT(ctx Node, _ string) string {
	var b strings.Builder
	b.WriteString("package ")
	b.WriteString(p.Name.String())
	b.WriteRune('\n')
	if len(p.Interfaces) > 0 {
		b.WriteRune('\n')
		for i, name := range codec.SortedKeys(p.Interfaces) {
			if i > 0 {
				b.WriteRune('\n')
			}
			b.WriteString(p.Interfaces[name].WIT(p, name))
			b.WriteRune('\n')
		}
	}
	if len(p.Worlds) > 0 {
		b.WriteRune('\n')
		for i, name := range codec.SortedKeys(p.Worlds) {
			if i > 0 {
				b.WriteRune('\n')
			}
			b.WriteString(p.Worlds[name].WIT(p, name))
			b.WriteRune('\n')
		}
	}
	return b.String()
}
