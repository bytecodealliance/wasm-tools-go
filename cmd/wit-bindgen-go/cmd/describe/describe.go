package describe

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/urfave/cli/v3"
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

	for _, w := range res.Worlds {
		describeWorld(p, w)
	}

	return nil
}

func describeWorld(p *printer, w *wit.World) {
	name := w.Package.Name.String() + "/" + w.Name
	if len(w.Imports) == 0 && len(w.Exports) == 0 {
		p.Printf("world %s {}", name)
		return
	}
	// TODO: print World.Docs
	p.Printf("world %s {", name)
	{
		p := p.indent()
		for name, item := range w.Imports {
			describeWorldItem(p, "import ", name, item)
		}
		for name, item := range w.Exports {
			describeWorldItem(p, "export ", name, item)
		}
	}
	p.Println("}")
	p.Println()
}

func describeWorldItem(p *printer, pfx, name string, item wit.WorldItem) {
	switch v := item.(type) {
	case *wit.Interface:
		describeWorldInterface(p, pfx, name, v)
	case *wit.TypeDef:
		describeWorldTypeDef(p, pfx, name, v)
	case *wit.Function:
		describeWorldFunction(p, pfx, name, v)
	}
}

func describeWorldInterface(p *printer, pfx, name string, i *wit.Interface) {
	if i.Name != nil {
		name = i.Package.Name.String() + "/" + *i.Name
	}
	// TODO: print Interface.Docs
	if len(i.TypeDefs) == 0 && len(i.Functions) == 0 {
		p.Printf("%s%s {}", pfx, name)
		return
	}
	p.Printf("%s%s {", pfx, name)
	{

	}
	p.Println("}")
}

func describeWorldTypeDef(p *printer, pfx, name string, t *wit.TypeDef) {
	// TODO
}

func describeWorldFunction(p *printer, pfx, name string, t *wit.Function) {
	// TODO
}

type printer struct {
	w     io.Writer
	depth int
}

func (p *printer) indent() *printer {
	pi := *p
	pi.depth++
	return &pi
}

func (p *printer) println(s string) {
	fmt.Fprintln(p.w, strings.Repeat("\t", p.depth), s)
}

func (p *printer) Println(a ...any) {
	s := fmt.Sprint(a...)
	for _, line := range strings.Split(s, "\n") {
		p.println(line)
	}
}

func (p *printer) Printf(format string, a ...any) {
	s := fmt.Sprintf(format, a...)
	for _, line := range strings.Split(s, "\n") {
		p.println(line)
	}
}
