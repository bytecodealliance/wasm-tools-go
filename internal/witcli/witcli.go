package witcli

import (
	"fmt"
	"strings"

	"github.com/bytecodealliance/wasm-tools-go/wit"
)

// LoadOne loads a single [wit.Resolve].
// An error is returned if len(paths) > 1.
// If paths is empty, or paths[0] == "" or "-", then it reads from stdin.
// If the resolved path doesn’t end in ".json", it will attempt to load
// WIT indirectly by processing the input through wasm-tools.
// If forceWIT is true, it will always process input through wasm-tools.
func LoadOne(forceWIT bool, paths ...string) (*wit.Resolve, error) {
	var path string
	switch len(paths) {
	case 0:
		path = "-"
	case 1:
		path = paths[0]
	default:
		return nil, fmt.Errorf("found %d path arguments, expecting 0 or 1", len(paths))
	}
	if forceWIT || !strings.HasSuffix(path, ".json") {
		return wit.LoadWIT(path)
	}
	return wit.LoadJSON(path)
}
