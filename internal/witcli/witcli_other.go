//go:build !wasip1 && !wasip2 && !tinygo

package witcli

import (
	"fmt"
	"io/fs"
	"os"
)

// FindOrCreateDir attempts to load file information for the provided path and
// the provided path does not already exist, it attempts to create the directory
// and then return the file information for it.
func FindOrCreateDir(path string) (fs.FileInfo, error) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = os.MkdirAll(path, os.FileMode(os.ModeDir|os.ModePerm))
		if err != nil {
			return nil, err
		}
		info, err = os.Stat(path)
	}
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", path)
	}
	return info, nil
}
