package describe

import (
	"fmt"

	"github.com/k0kubun/pp/v3"
	"github.com/urfave/cli/v3"
	"github.com/ydnar/wasm-tools-go/internal/witcli"
)

// Command is the CLI command for describe.
var Command = &cli.Command{
	Name:   "describe",
	Usage:  "describe a WIT JSON file",
	Action: action,
}

func action(ctx *cli.Context) error {
	res, err := witcli.LoadOne(ctx.Args().Slice()...)
	if err != nil {
		return err
	}

	fmt.Printf("// %d worlds(s), %d packages(s), %d interfaces(s), %d types(s)\n",
		len(res.Worlds), len(res.Packages), len(res.Interfaces), len(res.TypeDefs))
	p := pp.New()
	p.SetExportedOnly(true)
	p.Print(res)

	return nil
}
