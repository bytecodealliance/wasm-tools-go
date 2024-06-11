package wit

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// LoadJSON loads a [WIT] JSON file from path.
// If path is "" or "-", it reads from os.Stdin.
//
// [WIT]: https://github.com/WebAssembly/component-model/blob/main/design/mvp/WIT.md
func LoadJSON(path string) (*Resolve, error) {
	if path == "" || path == "-" {
		return DecodeJSON(os.Stdin)
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return DecodeJSON(f)
}

// LoadWIT loads [WIT] data from path by processing it through [wasm-tools].
// This will fail if wasm-tools is not in $PATH.
// If path is "" or "-", it reads from os.Stdin.
//
// [WIT]: https://github.com/WebAssembly/component-model/blob/main/design/mvp/WIT.md
// [wasm-tools]: https://crates.io/crates/wasm-tools
func LoadWIT(path string) (*Resolve, error) {
	wasmTools, err := exec.LookPath("wasm-tools")
	if err != nil {
		return nil, err
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command(wasmTools, "component", "wit", "-j", "--all-features")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if path == "" || path == "-" {
		cmd.Stdin = os.Stdin
	} else {
		cmd.Args = append(cmd.Args, path)
	}

	fmt.Printf("%s\n", strings.Join(cmd.Args, " "))

	err = cmd.Run()
	if err != nil {
		fmt.Fprint(os.Stderr, stderr.String())
		return nil, err
	}

	return DecodeJSON(&stdout)
}
