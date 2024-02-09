package generate

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
	"github.com/ydnar/wasm-tools-go/bindgen"
	"github.com/ydnar/wasm-tools-go/internal/codec"
	"github.com/ydnar/wasm-tools-go/internal/go/gen"
	"github.com/ydnar/wasm-tools-go/internal/witcli"
)

// Command is the CLI command for wit.
var Command = &cli.Command{
	Name:    "generate",
	Aliases: []string{"go"},
	Usage:   "generates Go from a fully-resolved WIT JSON file",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:      "out",
			Aliases:   []string{"o"},
			Value:     ".",
			TakesFile: true,
			OnlyOnce:  true,
			Usage:     "output directory",
		},
		&cli.StringMapFlag{
			Name:  "map",
			Usage: "maps WIT identifiers to Go identifiers",
		},
	},
	Action: action,
}

func action(ctx context.Context, cmd *cli.Command) error {
	out := cmd.String("out")
	if !isDir(out) {
		return fmt.Errorf("%s is not a directory", out)
	}
	fmt.Fprintf(os.Stderr, "Output dir: %s\n", out)

	pkgPath, err := gen.PackagePath(out)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Package root: %s\n", pkgPath)

	res, err := witcli.LoadOne(cmd.Args().Slice()...)
	if err != nil {
		return err
	}

	packages, err := bindgen.Go(res,
		bindgen.GeneratedBy(cmd.Root().Name),
		bindgen.PackageRoot(pkgPath),
		bindgen.MapIdents(cmd.StringMap("map")),
	)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Generated %d package(s)\n", len(packages))

	for _, pkg := range packages {
		for _, filename := range codec.SortedKeys(pkg.Files) {
			file := pkg.Files[filename]
			fmt.Fprintf(os.Stderr, "Generated file: %s/%s\n\n", pkg.Path, file.Name)
			b, err := file.Bytes()
			if err != nil {
				return err
			}
			fmt.Println(string(b))
		}
		fmt.Fprintln(os.Stderr)
	}

	return nil
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
