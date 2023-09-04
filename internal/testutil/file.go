package testutil

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
)

// ReadFile reads a file, returning the contents or an error.
func ReadFile(p string) ([]byte, error) {
	f, err := os.Open(Path(p))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
}

// Path returns an absolute path to the source-file relative path p.
// Used for testing when PWD is not set, such as running tests under wasm/wasip1.
func Path(p string) string {
	if !filepath.IsLocal(p) {
		return p
	}
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return p
	}
	dir := filepath.Dir(file)
	return filepath.Join(dir, p)
}
