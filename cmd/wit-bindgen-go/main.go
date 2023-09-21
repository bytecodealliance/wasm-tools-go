package main

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
	"github.com/ydnar/wasm-tools-go/cmd/wit-bindgen-go/cmd/describe"
)

func main() {
	cmd := &cli.Command{
		Name: "wit-bindgen-go",
		Commands: []*cli.Command{
			describe.Command,
		},
	}

	err := cmd.Run(context.Background(), os.Args)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}
