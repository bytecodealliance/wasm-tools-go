package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"slices"

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
			err := summarize(os.Stdin, "STDIN")
			if err != nil {
				return err
			}
		}
		f, err := os.Open(arg)
		if err != nil {
			return err
		}
		err = summarize(f, arg)
		f.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func summarize(r io.Reader, name string) error {
	res, err := wit.DecodeJSON(r)
	if err != nil {
		return err
	}

	fmt.Printf("Summarizing %s\n", name)
	fmt.Println()
	fmt.Printf("%d worlds(s), %d packages(s), %d interfaces(s), %d types(s)\n",
		len(res.Worlds), len(res.Packages), len(res.Interfaces), len(res.TypeDefs))
	fmt.Println()

	for _, p := range res.Packages {
		summarizePackage(p, "")
	}

	for i, t := range res.TypeDefs {
		fmt.Printf("Type %d: %s\n", i, t.Name)
		if t.Owner != nil {
			fmt.Printf("Owner: ")
			switch owner := t.Owner.(type) {
			case *wit.Interface:
				fmt.Printf("interface(%s)\n", Some(owner.Name, "<no name>"))
			case *wit.World:
				fmt.Printf("world(%s)\n", owner.Name)
			}
		}
		fmt.Println()
	}

	return nil
}

func summarizePackage(p *wit.Package, indent string) {
	fmt.Printf("Package: %s\n", string(p.Name))
	fmt.Printf("%d worlds, %d interface(s)\n", len(p.Worlds), len(p.Interfaces))
	fmt.Println()

	keys := Keys(p.Worlds)
	slices.Sort(keys)
	for _, k := range keys {
		summarizeWorld(p.Worlds[k], indent+"\t")
	}

	keys = Keys(p.Interfaces)
	slices.Sort(keys)
	for _, k := range keys {
		summarizeInterface(p.Interfaces[k], indent+"\t")
	}
}

func summarizeWorld(w *wit.World, indent string) {
	fmt.Printf("%sWorld: %s\n", indent, w.Name)
	fmt.Printf("%s%d import(s), %d export(s)\n", indent, len(w.Imports), len(w.Exports))
	fmt.Println()

	indent += "\t"

	keys := Keys(w.Imports)
	slices.Sort(keys)
	for _, k := range keys {
		fmt.Printf("%sImport: %s\n", indent, k)
	}
	fmt.Println()

	keys = Keys(w.Exports)
	slices.Sort(keys)
	for _, k := range keys {
		fmt.Printf("%sExport: %s\n", indent, k)
	}
	fmt.Println()
}

func summarizeInterface(i *wit.Interface, indent string) {
	fmt.Printf("%sInterface: %s\n", indent, Some(i.Name, "<unnamed>"))
	fmt.Printf("%s%d type(s), %d function(s)\n", indent, len(i.TypeDefs), len(i.Functions))
	fmt.Println()
}

func Some(s *string, fallback string) string {
	if s != nil {
		return *s
	}
	return fallback
}

func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
