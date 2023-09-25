package describe

import (
	"fmt"

	"github.com/urfave/cli/v3"
	"github.com/ydnar/wasm-tools-go/internal/witcli"
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
	fmt.Println(res.WIT(nil, ""))
	return nil
}
