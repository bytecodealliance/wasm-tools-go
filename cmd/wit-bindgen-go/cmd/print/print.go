package print

import (
	"fmt"

	"github.com/ydnar/wasm-tools-go/internal/witcli"

	"github.com/kr/pretty"
	"github.com/urfave/cli/v3"
)

// Command is the CLI command for describe.
var Command = &cli.Command{
	Name:   "print",
	Usage:  "pretty-prints a Resolve struct loaded from a WIT JSON file",
	Action: action,
}

func action(ctx *cli.Context) error {
	res, err := witcli.LoadOne(ctx.Args().Slice()...)
	if err != nil {
		return err
	}

	fmt.Printf("// %d worlds(s), %d packages(s), %d interfaces(s), %d types(s)\n",
		len(res.Worlds), len(res.Packages), len(res.Interfaces), len(res.TypeDefs))
	pretty.Print(res)

	return nil
}
