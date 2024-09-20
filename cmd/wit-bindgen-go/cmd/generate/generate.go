package generate

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
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

type Config struct {
	DryRun      bool
	OutDir      string
	OutPerm     os.FileMode
	PackageRoot string
	World       string
	CMPackage   string
	Versioned   bool
	ForceWIT    bool
	Path        string
}

func action(ctx context.Context, cmd *cli.Command) error {
	config, err := parseFlags(cmd)
	if err != nil {
		return err
	}

	res, err := loadWITModule(ctx, config)
	if err != nil {
		return err
	}

	packages, err := bindgen.Go(res,
		bindgen.GeneratedBy(cmd.Root().Name),
		bindgen.World(config.World),
		bindgen.PackageRoot(config.PackageRoot),
		bindgen.Versioned(config.Versioned),
		bindgen.CMPackage(config.CMPackage),
	)
	if err != nil {
		return err
	}

	if err := writeGoPackages(packages, config); err != nil {
		return err
	}

	return writeWITPackage(res, config)
}

func parseFlags(cmd *cli.Command) (*Config, error) {
	dryRun := cmd.Bool("dry-run")
	out := cmd.String("out")

	info, err := os.Stat(out)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", out)
	}
	fmt.Fprintf(os.Stderr, "Output dir: %s\n", out)
	outPerm := info.Mode().Perm()

	pkgRoot := cmd.String("package-root")
	if !cmd.IsSet("package-root") {
		pkgRoot, err = gen.PackagePath(out)
		if err != nil {
			return nil, err
		}
	}
	fmt.Fprintf(os.Stderr, "Package root: %s\n", pkgRoot)

	path, err := witcli.ParsePaths(cmd.Args().Slice()...)
	if err != nil {
		return nil, err
	}

	return &Config{
		DryRun:      dryRun,
		OutDir:      out,
		OutPerm:     outPerm,
		PackageRoot: pkgRoot,
		World:       cmd.String("world"),
		CMPackage:   cmd.String("cm"),
		Versioned:   cmd.Bool("versioned"),
		ForceWIT:    cmd.Bool("force-wit"),
		Path:        path,
	}, nil
}

func loadWITModule(ctx context.Context, config *Config) (*wit.Resolve, error) {
	if oci.IsOCIPath(config.Path) {
		fmt.Fprintf(os.Stderr, "Fetching OCI artifact %s\n", config.Path)
		bytes, err := oci.PullWIT(ctx, config.Path)
		if err != nil {
			return nil, err
		}
		return wit.LoadWITFromBuffer(bytes.Bytes())
	}

	return witcli.LoadOne(config.ForceWIT, config.Path)
}

func writeGoPackages(packages []*gen.Package, config *Config) error {
	fmt.Fprintf(os.Stderr, "Generated %d package(s)\n", len(packages))
	for _, pkg := range packages {
		if !pkg.HasContent() {
			fmt.Fprintf(os.Stderr, "Skipping empty package: %s\n", pkg.Path)
			continue
		}
		fmt.Fprintf(os.Stderr, "Generated package: %s\n", pkg.Path)

		for _, filename := range codec.SortedKeys(pkg.Files) {
			file := pkg.Files[filename]
			dir := filepath.Join(config.OutDir, strings.TrimPrefix(file.Package.Path, config.PackageRoot))
			path := filepath.Join(dir, file.Name)

			if !file.HasContent() {
				fmt.Fprintf(os.Stderr, "Skipping empty file: %s\n", path)
				continue
			}

			if err := os.MkdirAll(dir, config.OutPerm); err != nil {
				return err
			}

			content, err := file.Bytes()
			if err != nil {
				if content == nil {
					return err
				}
				fmt.Fprintf(os.Stderr, "Error formatting file: %v\n", err)
			} else {
				fmt.Fprintf(os.Stderr, "Generated file: %s\n", path)
			}

			if config.DryRun {
				fmt.Println(string(content))
				fmt.Println()
				continue
			}

			if err := os.WriteFile(path, content, config.OutPerm); err != nil {
				return err
			}
		}
	}
	return nil
}

func writeWITPackage(res *wit.Resolve, config *Config) error {
	witDir := filepath.Join(config.OutDir, "wit")
	if err := os.MkdirAll(witDir, config.OutPerm); err != nil {
		return err
	}
	witFilePath := filepath.Join(witDir, "webassembly.wit.wasm")
	fmt.Fprintf(os.Stderr, "Generated WIT file: %s\n", witFilePath)

	wasmTools, err := exec.LookPath("wasm-tools")
	if err != nil {
		return err
	}

	var stderr bytes.Buffer
	wasmCmd := exec.Command(wasmTools, "component", "wit", "--wasm", "--all-features", "--output", witFilePath)
	wasmCmd.Stderr = &stderr
	wasmCmd.Stdin = bytes.NewReader([]byte(res.WIT(nil, "")))

	if err := wasmCmd.Run(); err != nil {
		fmt.Fprint(os.Stderr, stderr.String())
		return err
	}
	return nil
}
