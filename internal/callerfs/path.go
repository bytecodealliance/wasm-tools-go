package callerfs

import (
	"path/filepath"
	"runtime"
)

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
