package wit

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/ydnar/wasm-tools-go/internal/callerfs"
)

const testdataDir = "../testdata/"

func loadTestdata(f func(path string, res *Resolve) error) error {
	return filepath.WalkDir(callerfs.Path(testdataDir), func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fs.SkipDir
		}
		if !strings.HasSuffix(path, ".wit.json") && !strings.HasSuffix(path, ".wit.md.json") {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		res, err := DecodeJSON(file)
		if err != nil {
			return err
		}
		return f(path, res)
	})
}
