package testdata

import (
	"io/fs"
	"path/filepath"

	"github.com/ydnar/wasm-tools-go/internal/callerfs"
)

// Walk walks the files in the testdata directory, passing them to func.
// Supply glob patterns (e.g. "*.wit.json") to filter files passed to f.
func Walk(f func(path string) error, patterns ...string) error {
	return filepath.WalkDir(callerfs.Path("."), func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fs.SkipDir
		}
		if len(patterns) == 0 {
			return f(path)
		}
		for _, p := range patterns {
			matched, err := filepath.Match(p, filepath.Base(path))
			if err != nil {
				return err
			}
			if matched {
				return f(path)
			}
		}
		return nil
	})
}

// Relative returns a testdata-relative path.
// If path is not relative to the testdata directory, path is returned unchanged.
func Relative(path string) string {
	rel, err := filepath.Rel(callerfs.Path("."), path)
	if err != nil {
		return path
	}
	return rel
}
