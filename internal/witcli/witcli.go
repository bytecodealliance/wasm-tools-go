package witcli

import (
	"fmt"
	"strings"

	"github.com/bytecodealliance/wasm-tools-go/wit"
)

// LoadOne loads a single [wit.Resolve].
// If path == "" or "-", then it reads from stdin.
// If the resolved path doesn’t end in ".json", it will attempt to load
// WIT indirectly by processing the input through wasm-tools.
// If forceWIT is true, it will always process input through wasm-tools.
func LoadOne(forceWIT bool, path string) (*wit.Resolve, error) {
	if forceWIT || !strings.HasSuffix(path, ".json") {
		return wit.LoadWITFromPath(path)
	}
	return wit.LoadJSON(path)
}

// An error is returned if len(paths) > 1.
// If paths is empty, returns "-".
func ParsePaths(paths ...string) (string, error) {
	var path string
	switch len(paths) {
	case 0:
		path = "-"
	case 1:
		path = paths[0]
	default:
		return "", fmt.Errorf("found %d path arguments, expecting 0 or 1", len(paths))
	}
	return path, nil
}
