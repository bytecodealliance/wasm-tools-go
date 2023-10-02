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

func (_node) WIT(ctx Node, name string) string { return "/* TODO(" + name + ") */" }

func indent(s string) string {
	const ws = "    "
	return strings.TrimSuffix(ws+strings.ReplaceAll(s, "\n", "\n"+ws), ws)
}

// unwrap unwraps a multiline string into a single line, if:
// 1. its length is <= 50 chars
// 2. its line count is <= 5
// This is used for single-line [Record], [Flags], [Variant], and [Enum] declarations.
func unwrap(s string) string {
	const chars = 50
	const lines = 5
	if len(s) > chars || strings.Count(s, "\n") > lines {
		return s
	}
	var b strings.Builder
	for i, line := range strings.Split(s, "\n") {
		if i > 0 {
			b.WriteRune(' ')
		}
		b.WriteString(strings.Trim(line, " \t\r\n"))
	}
	return b.String()
}

// WIT returns the WIT representation of r.
func (r *Resolve) WIT(_ Node, _ string) string {
	var b strings.Builder
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
	if name == "" {
		name = w.Name
	}
	var b strings.Builder
	// TODO: docs
	fmt.Fprintf(&b, "world %s {", name) // TODO: compare to w.Name?
	if len(w.Imports) > 0 || len(w.Exports) > 0 {
		b.WriteRune('\n')
		for _, name := range codec.SortedKeys(w.Imports) {
			b.WriteString(indent(w.itemWIT("import", name, w.Imports[name])))
			b.WriteString(";\n")
		}
		for _, name := range codec.SortedKeys(w.Exports) {
			b.WriteString(indent(w.itemWIT("export", name, w.Exports[name])))
			b.WriteString(";\n")
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
		rname := relativeName(i, ctx.Package)
		if rname != "" {
			return rname
		}

		// Otherwise, this is an inline interface decl.
		b.WriteString(name)
		b.WriteString(": interface ")
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
			b.WriteString(";\n")
			n++
		}
		for _, name := range codec.SortedKeys(i.Functions) {
			// if n > 0 {
			// 	b.WriteRune('\n')
			// }
			b.WriteString(indent(i.Functions[name].WIT(i, name)))
			b.WriteString(";\n")
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
		ownerName := relativeName(t.Owner, ctx.Package())
		if t.Name != nil && *t.Name != name {
			return fmt.Sprintf("use %s.{%s as %s}", ownerName, *t.Name, name)
		}
		return fmt.Sprintf("use %s.{%s}", ownerName, name)

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

func relativeName(o TypeOwner, p *Package) string {
	var op *Package
	var name string
	switch o := o.(type) {
	case *Interface:
		if o.Name == nil {
			return ""
		}
		op = o.Package
		name = *o.Name

	case *World:
		op = o.Package
		name = o.Name
	}
	if op == p {
		return name
	}
	if op == nil {
		return ""
	}
	return op.Name.String() + "/" + name
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
	b.WriteRune('}')
	return unwrap(b.String())
}

func (f *Field) WIT(ctx Node, name string) string {
	// TODO: docs
	return f.Name + ": " + f.Type.WIT(f, "")
}

func (r *Resource) WIT(ctx Node, name string) string {
	var b strings.Builder
	b.WriteString("resource ")
	b.WriteString(name)
	b.WriteString(" {} // TODO: constructor, methods, and static functions")
	return b.String()
}

func (h *OwnedHandle) WIT(ctx Node, name string) string {
	var b strings.Builder
	if name != "" {
		b.WriteString("type ")
		b.WriteString(name)
		b.WriteString(" = ")
	}
	b.WriteString("own<")
	b.WriteString(h.Type.WIT(h, ""))
	b.WriteRune('>')
	return b.String()
}

func (h *BorrowedHandle) WIT(ctx Node, name string) string {
	var b strings.Builder
	if name != "" {
		b.WriteString("type ")
		b.WriteString(name)
		b.WriteString(" = ")
	}
	b.WriteString("borrow<")
	b.WriteString(h.Type.WIT(h, ""))
	b.WriteRune('>')
	return b.String()
}

func (f *Flags) WIT(ctx Node, name string) string {
	var b strings.Builder
	b.WriteString("flags ")
	b.WriteString(name)
	b.WriteString(" {")
	if len(f.Flags) > 0 {
		for i := range f.Flags {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(f.Flags[i].WIT(f, ""))
		}
	}
	b.WriteRune('}')
	return unwrap(b.String())
}

func (f *Flag) WIT(_ Node, _ string) string {
	// TODO: docs
	return f.Name
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

func (v *Variant) WIT(_ Node, name string) string {
	var b strings.Builder
	b.WriteString("variant ")
	b.WriteString(name)
	b.WriteString(" {")
	if len(v.Cases) > 0 {
		b.WriteRune('\n')
		for i := range v.Cases {
			if i > 0 {
				b.WriteString(",\n")
			}
			b.WriteString(indent(v.Cases[i].WIT(v, "")))
		}
		b.WriteRune('\n')
	}
	b.WriteRune('}')
	return unwrap(b.String())
}

func (c *Case) WIT(_ Node, _ string) string {
	// TODO: docs
	var b strings.Builder
	b.WriteString(c.Name)
	if c.Type != nil {
		b.WriteRune('(')
		b.WriteString(c.Type.WIT(c, ""))
		b.WriteRune(')')
	}
	return b.String()
}

func (e *Enum) WIT(_ Node, name string) string {
	var b strings.Builder
	b.WriteString("enum ")
	b.WriteString(name)
	b.WriteString(" {")
	if len(e.Cases) > 0 {
		b.WriteRune('\n')
		for i := range e.Cases {
			if i > 0 {
				b.WriteString(",\n")
			}
			b.WriteString(indent(e.Cases[i].WIT(e, "")))
		}
		b.WriteRune('\n')
	}
	b.WriteRune('}')
	return unwrap(b.String())
}

func (c *EnumCase) WIT(_ Node, _ string) string {
	// TODO: docs
	return c.Name
}

func (o *Option) WIT(_ Node, name string) string {
	var b strings.Builder
	if name != "" {
		b.WriteString("type ")
		b.WriteString(name)
		b.WriteString(" = ")
	}
	b.WriteString("option<")
	b.WriteString(o.Type.WIT(o, ""))
	b.WriteRune('>')
	return b.String()
}

func (r *Result) WIT(_ Node, name string) string {
	var b strings.Builder
	if name != "" {
		b.WriteString("type ")
		b.WriteString(name)
		b.WriteString(" = ")
	}
	b.WriteString("result<")
	if r.OK != nil {
		b.WriteString(r.OK.WIT(r, ""))
		b.WriteString(", ")
	} else {
		b.WriteString("_, ")
	}
	if r.Err != nil {
		b.WriteString(r.Err.WIT(r, ""))
	} else {
		b.WriteRune('_')
	}
	b.WriteRune('>')
	return b.String()
}

func (l *List) WIT(_ Node, name string) string {
	var b strings.Builder
	if name != "" {
		b.WriteString("type ")
		b.WriteString(name)
		b.WriteString(" = ")
	}
	b.WriteString("list<")
	b.WriteString(l.Type.WIT(l, ""))
	b.WriteRune('>')
	return b.String()
}

func (f *Future) WIT(_ Node, name string) string {
	var b strings.Builder
	if name != "" {
		b.WriteString("type ")
		b.WriteString(name)
		b.WriteString(" = ")
	}
	b.WriteString("future")
	if f.Type != nil {
		b.WriteRune('<')
		b.WriteString(f.Type.WIT(f, ""))
		b.WriteRune('>')
	}
	return b.String()
}

func (s *Stream) WIT(_ Node, name string) string {
	var b strings.Builder
	if name != "" {
		b.WriteString("type ")
		b.WriteString(name)
		b.WriteString(" = ")
	}
	b.WriteString("stream<")
	if s.Element != nil {
		b.WriteString(s.Element.WIT(s, ""))
		b.WriteString(", ")
	} else {
		b.WriteString("_, ")
	}
	if s.End != nil {
		b.WriteString(s.End.WIT(s, ""))
	} else {
		b.WriteRune('_')
	}
	b.WriteRune('>')
	return b.String()
}

// WIT returns the WIT representation of [primitive type] T.
//
// [primitive type]: https://component-model.bytecodealliance.org/wit-overview.html#primitive-types
func (p _primitive[T]) WIT(_ Node, name string) string {
	if name != "" {
		return "type " + name + " = " + p.String()
	}
	return p.String()
}

// WIT returns the WIT representation of f.
func (f *Function) WIT(_ Node, name string) string {
	// TODO: docs
	var b strings.Builder
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
	// TODO: docs
	var b strings.Builder
	b.WriteString("package ")
	b.WriteString(p.Name.String())
	b.WriteString(";\n")
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
