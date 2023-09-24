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

	p := &printer{w: os.Stdout}

	for i, w := range res.Worlds {
		if i > 0 {
			p.Println()
		}
		describeWorld(p, w)
	}

	return nil
}

func describeWorld(p *printer, w *wit.World) {
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
		for i, name := range codec.SortedKeys(w.Imports) {
			if i > 0 {
				p.Println()
			}
			p.Print("import ")
			describeWorldItem(p, name, w.Imports[name])
		}
		for i, name := range codec.SortedKeys(w.Exports) {
			if i > 0 {
				p.Println()
			}
			p.Print("export ")
			describeWorldItem(p, name, w.Exports[name])
		}
	}
	p.Println("}")
}

func describeWorldItem(p *printer, name string, item wit.WorldItem) {
	switch v := item.(type) {
	case *wit.Interface:
		describeInterface(p, name, v)
	case *wit.TypeDef:
		describeTypeDef(p, name, v)
	case *wit.Function:
		describeFunction(p, name, v)
	}
}

func describeInterface(p *printer, name string, iface *wit.Interface) {
	if iface.Name != nil {
		name = iface.Package.Name.String() + "/" + *iface.Name
	}
	// TODO: print Interface.Docs
	p.Printf("%s {", name)
	if len(iface.TypeDefs) > 0 || len(iface.Functions) > 0 {
		p.Println()
		p := p.Indent()
		for i, name := range codec.SortedKeys(iface.TypeDefs) {
			if i > 0 {
				fmt.Println()
			}
			describeTypeDef(p, name, iface.TypeDefs[name])
		}
		if len(iface.TypeDefs) > 0 && len(iface.Functions) > 0 {
			fmt.Println()
		}
		for i, name := range codec.SortedKeys(iface.Functions) {
			if i > 0 {
				fmt.Println()
			}
			describeFunction(p, name, iface.Functions[name])
		}
	}
	p.Println("}")
}

func describeTypeDef(p *printer, name string, t *wit.TypeDef) {
	p.Printf("type %s = ", name)
	p.Println("TypeDef(TODO)")
}

func describeFunction(p *printer, name string, f *wit.Function) {
	// TODO: print Function.Docs
	p.Printf("%s: func(", name)
	describeParams(p, f.Params)
	p.Printf(")")
	if len(f.Results) > 0 {
		p.Printf(" -> ")
		describeParams(p, f.Results)
	}
	p.Println()
}

func describeParams(p *printer, params []wit.Param) {
	for i, param := range params {
		if i > 0 {
			p.Print(", ")
		}
		if param.Name != "" {
			p.Printf("%s: ", param.Name)
		}
		describeType(p, param.Type)
	}
}

func describeType(p *printer, t wit.Type) {
	p.Print("T")
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
