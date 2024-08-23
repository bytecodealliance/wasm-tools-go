package main

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/urfave/cli/v3"

	"github.com/ydnar/wasm-tools-go/cmd/wit-bindgen-go/cmd/generate"
	"github.com/ydnar/wasm-tools-go/cmd/wit-bindgen-go/cmd/wit"
)

var (
	version  = ""
	revision = ""
)

func init() {
	build, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}
	version = build.Main.Version
	for _, s := range build.Settings {
		switch s.Key {
		case "vcs.revision":
			revision = s.Value
		}
	}
}

func main() {
	cmd := &cli.Command{
		Name:  "wit-bindgen-go",
		Usage: "inspect or manipulate WebAssembly Interface Types for Go",
		Commands: []*cli.Command{
			generate.Command,
			wit.Command,
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "force-wit",
				Usage: "force loading WIT via wasm-tools",
			},
		},
		Version: fmt.Sprintf("%v (%v)", version, revision),
	}

	err := cmd.Run(context.Background(), os.Args)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}
