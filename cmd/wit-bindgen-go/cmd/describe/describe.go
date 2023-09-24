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
		describeWorldInterface(p, name, v)
	case *wit.TypeDef:
		describeWorldTypeDef(p, name, v)
	case *wit.Function:
		describeWorldFunction(p, name, v)
	}
}

func describeWorldInterface(p *printer, name string, i *wit.Interface) {
	if i.Name != nil {
		name = i.Package.Name.String() + "/" + *i.Name
	}
	// TODO: print Interface.Docs
	p.Printf("%s {", name)
	if len(i.TypeDefs) > 0 || len(i.Functions) > 0 {
		p.Println()
		p := p.Indent()
		var _ = p
	}
	p.Println("}")
}

func describeWorldTypeDef(p *printer, name string, t *wit.TypeDef) {
	// TODO
}

func describeWorldFunction(p *printer, name string, t *wit.Function) {
	// TODO
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
