package describe

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/urfave/cli/v3"
	"github.com/ydnar/wasm-tools-go/internal/codec"
	"github.com/ydnar/wasm-tools-go/internal/witcli"
	"github.com/ydnar/wasm-tools-go/wit"
)

// Command is the CLI command for describe.
var Command = &cli.Command{
	Name:   "describe",
	Usage:  "describes a WIT JSON file",
	Action: action,
}

func action(ctx *cli.Context) error {
	res, err := witcli.LoadOne(ctx.Args().Slice()...)
	if err != nil {
		return err
	}

	fmt.Println(res.WIT())
	return nil

	p := &printer{w: os.Stdout}

	for i, w := range res.Worlds {
		if i > 0 {
			p.Println()
		}
		printWorld(p, w)
	}

	return nil
}

func printWorld(p *printer, w *wit.World) {
	name := w.Package.Name.String() + "/" + w.Name
	if len(w.Imports) == 0 && len(w.Exports) == 0 {
		p.Printf("world %s {}\n", name)
		return
	}
	// TODO: print World.Docs
	p.Printf("world %s {", name)
	if len(w.Imports) > 0 || len(w.Exports) > 0 {
		p.Println()
		p := p.Indent()
		n := 0
		for _, name := range codec.SortedKeys(w.Imports) {
			if n > 0 {
				p.Println()
			}
			p.Print("import ")
			printWorldItem(p, name, w.Imports[name])
			n++
		}
		for _, name := range codec.SortedKeys(w.Exports) {
			if n > 0 {
				p.Println()
			}
			p.Print("export ")
			printWorldItem(p, name, w.Exports[name])
			n++
		}
	}
	p.Println("}")
}

func printWorldItem(p *printer, name string, item wit.WorldItem) {
	switch v := item.(type) {
	case *wit.Interface:
		printInterface(p, name, v)
	case *wit.TypeDef:
		printTypeDef(p, name, v)
	case *wit.Function:
		printFunction(p, name, v)
	}
}

func printInterface(p *printer, name string, i *wit.Interface) {
	if i.Name != nil {
		name = i.Package.Name.String() + "/" + *i.Name
	}
	// TODO: print Interface.Docs
	p.Printf("%s {", name)
	if len(i.TypeDefs) > 0 || len(i.Functions) > 0 {
		p.Println()
		p := p.Indent()
		n := 0
		for _, name := range codec.SortedKeys(i.TypeDefs) {
			if n > 0 {
				fmt.Println()
			}
			printTypeDef(p, name, i.TypeDefs[name])
			n++
		}
		for _, name := range codec.SortedKeys(i.Functions) {
			if n > 0 {
				fmt.Println()
			}
			printFunction(p, name, i.Functions[name])
			n++
		}
	}
	p.Println("}")
}

func printTypeDef(p *printer, name string, t *wit.TypeDef) {
	if t.Name != nil {
		name = *t.Name // TODO: can we figure out if a type is imported from elsewhere?
	}
	p.Printf("type %s = ", name)
	p.Println("TypeDef(TODO)")
}

func printType(p *printer, t wit.Type) {
	switch t := t.(type) {
	case *wit.TypeDef:
		if t.Name != nil {
			p.Printf("%s", *t.Name)
			return
		}
	case interface{ WIT() string }:
		p.Printf("%s", t.WIT())
		return
	}
	p.Print("T")
}

func printFunction(p *printer, name string, f *wit.Function) {
	// TODO: print Function.Docs
	p.Printf("%s: func(", name)
	printParams(p, f.Params)
	p.Printf(")")
	if len(f.Results) > 0 {
		p.Printf(" -> ")
		printParams(p, f.Results)
	}
	p.Println()
}

func printParams(p *printer, params []wit.Param) {
	for i, param := range params {
		if i > 0 {
			p.Print(", ")
		}
		if param.Name != "" {
			p.Printf("%s: ", param.Name)
		}
		printType(p, param.Type)
	}
}

type printer struct {
	w        io.Writer
	depth    int
	indented int
}

func (p *printer) Indent() *printer {
	pi := *p
	pi.depth++
	return &pi
}

func (p *printer) Print(a ...any) {
	p.print(fmt.Sprint(a...))
}

func (p *printer) Println(a ...any) {
	p.print(fmt.Sprintln(a...))
}

func (p *printer) Printf(format string, a ...any) {
	p.print(fmt.Sprintf(format, a...))
}

func (p *printer) print(s string) {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		if i == len(lines)-1 && line == "" {
			return
		}
		fmt.Fprint(p.w, strings.Repeat("\t", p.depth-p.indented))
		p.indented = p.depth
		fmt.Fprint(p.w, line)
		if i < len(lines)-1 {
			fmt.Fprint(p.w, "\n")
			p.indented = 0
		}
	}
}
