package wit

import (
	"fmt"
	"strings"

	"github.com/ydnar/wasm-tools-go/internal/codec"
)

const witIndent = "    "

func indent(s string) string {
	return strings.TrimSuffix(witIndent+strings.Replace(s, "\n", "\n"+witIndent, 0), witIndent)
}

// Syntax represents any type that can return a WIT representation of itself.
type Syntax interface {
	WIT() string
}

// WIT returns the WIT representation of r.
func (r *Resolve) WIT() string {
	b := &strings.Builder{}
	for i, p := range r.Packages {
		if i > 0 {
			b.WriteRune('\n')
			b.WriteRune('\n')
		}
		fmt.Fprintf(b, "package %s\n", p.Name.String())
		if len(p.Interfaces) > 0 {
			b.WriteRune('\n')
			for i, name := range codec.SortedKeys(p.Interfaces) {
				if i > 0 {
					b.WriteRune('\n')
				}
				b.WriteString(fmt.Sprintf("interface %s %s\n", name, p.Interfaces[name].WIT()))
			}
		}
		if len(p.Worlds) > 0 {
			b.WriteRune('\n')
			for i, name := range codec.SortedKeys(p.Worlds) {
				if i > 0 {
					b.WriteRune('\n')
				}
				b.WriteString(fmt.Sprintf("world %s %s\n", name, p.Worlds[name].WIT()))
			}
		}
	}
	return b.String()
}

// WIT returns the WIT representation of w.
func (w *World) WIT() string {
	var b strings.Builder
	// TODO: docs
	b.WriteRune('{')
	if len(w.Imports) > 0 || len(w.Exports) > 0 {
		b.WriteRune('\n')
		for _, name := range codec.SortedKeys(w.Imports) {
			b.WriteString(indent(fmt.Sprintf("import %s\n", worldItemWIT(name, w.Imports[name]))))
		}
		for _, name := range codec.SortedKeys(w.Exports) {
			b.WriteString(indent(fmt.Sprintf("export %s\n", worldItemWIT(name, w.Exports[name]))))
		}
	}
	b.WriteRune('}')
	return b.String()
}

func worldItemWIT(name string, v WorldItem) string {
	switch v := v.(type) {
	case *Interface:
		if v.Name != nil {
			return fmt.Sprintf("%s/%s", v.Package.Name.String(), *v.Name)
		}
		return fmt.Sprintf("%s %s", name, v.WIT())
	case *TypeDef:
		switch owner := v.Owner.(type) {
		case *World:
			if v.Name != nil {
				// TODO: is this right?
				name = fmt.Sprintf("%s (%s???)", name, *v.Name)
			}
		case *Interface:
			if v.Name != nil {
				// TODO: is this right?
				return fmt.Sprintf("%s/%s", owner.Package.Name.String(), *v.Name)
			}
		}
		return fmt.Sprintf("%s %s", name, v.WIT())
	case *Function:
		return fmt.Sprintf("%s: %s", name, v.WIT())
	}
	panic("BUG: unknown WorldItem")
}

// WIT returns the WIT representation of i.
func (i *Interface) WIT() string {
	var b strings.Builder
	// TODO: docs
	b.WriteRune('{')
	if len(i.TypeDefs) > 0 || len(i.Functions) > 0 {
		b.WriteRune('\n')
		n := 0
		for _, name := range codec.SortedKeys(i.TypeDefs) {
			if n > 0 {
				b.WriteRune('\n')
			}
			b.WriteString(indent(fmt.Sprintf("type %s = %s\n", name, i.TypeDefs[name].WIT())))
			n++
		}
		for _, name := range codec.SortedKeys(i.Functions) {
			if n > 0 {
				b.WriteRune('\n')
			}
			b.WriteString(indent(fmt.Sprintf("%s: %s\n", name, i.Functions[name].WIT())))
			n++
		}
	}
	b.WriteRune('}')
	return b.String()
}

// WIT returns the WIT representation of t.
func (t *TypeDef) WIT() string {
	return "TODO<TypeDef>"
}

func (_type) WIT() string {
	return "TODO<_type>"
}

// WIT returns the WIT representation of f.
func (f *Function) WIT() string {
	var b strings.Builder
	// TODO: docs
	b.WriteString("func(")
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
		if param.Name != "" {
			b.WriteString(param.Name)
			b.WriteString(": ")
		}
		b.WriteString(param.Type.WIT())
	}
	return b.String()
}
