package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/ydnar/wasm-tools-go/wit"
)

func main() {
	err := Main()
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}
}

func Main() error {
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		args = []string{"-"}
	}
	for _, arg := range args {
		if arg == "-" {
			err := process(os.Stdin, "STDIN")
			if err != nil {
				return err
			}
		}
		f, err := os.Open(arg)
		if err != nil {
			return err
		}
		err = process(f, arg)
		f.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func process(r io.Reader, name string) error {
	res, err := wit.DecodeJSON(r)
	if err != nil {
		return err
	}

	fmt.Printf("WIT:        %s\n", name)
	fmt.Printf("Worlds:     %d\n", len(res.Worlds))
	fmt.Printf("Packages:   %d\n", len(res.Packages))
	fmt.Printf("Interfaces: %d\n", len(res.Interfaces))
	fmt.Printf("TypeDefs:   %d\n", len(res.TypeDefs))

	return nil
}
