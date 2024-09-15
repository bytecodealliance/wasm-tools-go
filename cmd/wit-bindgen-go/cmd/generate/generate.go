package generate

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bytecodealliance/wasm-tools-go/internal/codec"
	"github.com/bytecodealliance/wasm-tools-go/internal/go/gen"
	"github.com/bytecodealliance/wasm-tools-go/internal/oci"
	"github.com/bytecodealliance/wasm-tools-go/internal/witcli"
	"github.com/bytecodealliance/wasm-tools-go/wit"
	"github.com/bytecodealliance/wasm-tools-go/wit/bindgen"
	"github.com/urfave/cli/v3"
)

// Command is the CLI command for wit.
var Command = &cli.Command{
	Name:    "generate",
	Aliases: []string{"go"},
	Usage:   "generate Go bindings from from WIT (WebAssembly Interface Types)",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "world",
			Aliases:  []string{"w"},
			Value:    "",
			OnlyOnce: true,
			Config:   cli.StringConfig{TrimSpace: true},
			Usage:    "WIT world to generate, otherwise generate all worlds",
		},
		&cli.StringFlag{
			Name:      "out",
			Aliases:   []string{"o"},
			Value:     ".",
			TakesFile: true,
			OnlyOnce:  true,
			Config:    cli.StringConfig{TrimSpace: true},
			Usage:     "output directory",
		},
		&cli.StringFlag{
			Name:     "package-root",
			Aliases:  []string{"p"},
			Value:    "",
			OnlyOnce: true,
			Config:   cli.StringConfig{TrimSpace: true},
			Usage:    "Go package root, e.g. github.com/org/repo/internal",
		},
		&cli.StringFlag{
			Name:     "cm",
			Value:    "",
			OnlyOnce: true,
			Config:   cli.StringConfig{TrimSpace: true},
			Usage:    "Import path for the Component Model utility package, e.g. github.com/bytecodealliance/wasm-tools-go/cm",
		},
		&cli.BoolFlag{
			Name:  "versioned",
			Usage: "emit versioned Go package(s) for each WIT version",
		},
		&cli.BoolFlag{
			Name:  "dry-run",
			Usage: "do not write files; print to stdout",
		},
	},
	Action: action,
}

func action(ctx context.Context, cmd *cli.Command) error {
	dryRun := cmd.Bool("dry-run")

	out := cmd.String("out")
	info, err := os.Stat(out)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", out)
	}
	fmt.Fprintf(os.Stderr, "Output dir: %s\n", out)
	outPerm := info.Mode().Perm()

	pkgRoot := cmd.String("package-root")
	if !cmd.IsSet("package-root") {
		pkgRoot, err = gen.PackagePath(out)
		if err != nil {
			return err
		}
	}
	fmt.Fprintf(os.Stderr, "Package root: %s\n", pkgRoot)

	path, err := witcli.ParsePaths(cmd.Args().Slice()...)
	if err != nil {
		return err
	}

	// check if the path is a OCI path
	var res *wit.Resolve
	if oci.IsOCIPath(path) {
		fmt.Fprintf(os.Stderr, "Fetching OCI artifact %s\n", path)
		bytes, err := oci.PullWIT(ctx, path)
		if err != nil {
			return err
		}
		res, err = wit.DecodeJSON(bytes)
	} else {
		res, err = witcli.LoadOne(cmd.Bool("force-wit"), path)
	}

	if err != nil {
		return err
	}

	packages, err := bindgen.Go(res,
		bindgen.GeneratedBy(cmd.Root().Name),
		bindgen.World(cmd.String("world")),
		bindgen.PackageRoot(pkgRoot),
		bindgen.Versioned(cmd.Bool("versioned")),
		bindgen.CMPackage(cmd.String("cm")),
	)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "Generated %d package(s)\n", len(packages))

	for _, pkg := range packages {
		if !pkg.HasContent() {
			fmt.Fprintf(os.Stderr, "Skipping empty package: %s\n", pkg.Path)
			continue
		}

		fmt.Fprintf(os.Stderr, "Generated package: %s\n", pkg.Path)

		for _, filename := range codec.SortedKeys(pkg.Files) {
			file := pkg.Files[filename]

			dir := filepath.Join(out, strings.TrimPrefix(file.Package.Path, pkgRoot))
			path := filepath.Join(dir, file.Name)

			if !file.HasContent() {
				fmt.Fprintf(os.Stderr, "Skipping empty file: %s\n", path)
				continue
			}

			err := os.MkdirAll(dir, outPerm)
			if err != nil {
				return err
			}

			b, err := file.Bytes()
			if err != nil {
				if b == nil {
					return err
				}
				fmt.Fprintf(os.Stderr, "Error formatting file: %v\n", err)
			} else {
				fmt.Fprintf(os.Stderr, "Generated file: %s\n", path)
			}

			if dryRun {
				fmt.Println(string(b))
				fmt.Println()
				continue
			}

			f, err := os.Create(path)
			if err != nil {
				return err
			}
			n, err := f.Write(b)
			f.Close()
			if err != nil {
				return err
			}
			if n != len(b) {
				return fmt.Errorf("wrote %d bytes to %s, expected %d", n, path, len(b))
			}
		}
	}

	return nil
}
