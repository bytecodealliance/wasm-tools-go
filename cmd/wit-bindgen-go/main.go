package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/k0kubun/pp/v3"
	"github.com/ydnar/wasm-tools-go/wit"
)

func main() {
	err := Main()
	if err != nil {
		fmt.Printf("error: %v\n", err)
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
			err := describe(os.Stdin, "STDIN")
			if err != nil {
				return err
			}
		}
		f, err := os.Open(arg)
		if err != nil {
			return err
		}
		err = describe(f, arg)
		f.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func describe(r io.Reader, name string) error {
	res, err := wit.DecodeJSON(r)
	if err != nil {
		return err
	}

	fmt.Printf("// Describing WIT from %s\n", name)
	fmt.Printf("// %d worlds(s), %d packages(s), %d interfaces(s), %d types(s)\n",
		len(res.Worlds), len(res.Packages), len(res.Interfaces), len(res.TypeDefs))
	fmt.Println()

	p := pp.New()
	p.SetExportedOnly(true)
	p.Print(res)

	return nil
}
