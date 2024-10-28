//go:build wasip1 || wasip2 || tinygo

package witcli

import (
	"errors"
	"io/fs"
)

func FindOrCreateDir(path string) (fs.FileInfo, error) {
	return nil, errors.New("FindOrCreateDir not supported on WASI or TinyGo")
}
