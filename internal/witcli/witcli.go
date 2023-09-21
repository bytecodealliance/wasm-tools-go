package witcli

import (
	"fmt"
	"os"

	"github.com/ydnar/wasm-tools-go/wit"
)

// LoadOne loads one WIT JSON file from paths.
// An error is returned if len(paths) > 1.
// If paths is empty, or paths[0] == "" or "-", then it reads from stdin.
func LoadOne(paths ...string) (*wit.Resolve, error) {
	var path string
	switch len(paths) {
	case 0:
		path = "-"
	case 1:
		path = paths[0]
	default:
		return nil, fmt.Errorf("found %d path arguments, expecting 0 or 1", len(paths))
	}
	return LoadJSON(path)
}

// LoadJSON loads a WIT JSON file from path.
// If path is "" or "-", it reads from os.Stdin.
func LoadJSON(path string) (*wit.Resolve, error) {
	if path == "" || path == "-" {
		return wit.DecodeJSON(os.Stdin)
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return wit.DecodeJSON(f)
}
