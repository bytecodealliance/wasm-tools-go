package relpath

import (
	"os"
	"path/filepath"
)

// Abs returns a best-effort absolute representation of path.
// The returned path may be relative on platforms such as WebAssembly.
// If the path is not absolute it will be joined with the current
// working directory to turn it into an absolute path.
// See [filepath.Abs] for more information.
func Abs(path string) (string, error) {
	if filepath.IsAbs(path) {
		return filepath.Clean(path), nil
	}
	wd, err := Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(wd, path), nil
}

// Getwd returns best-effort path to the current working directory, even on WebAssembly.
// The path may be absolute, or it may be relative, prefixed with one or more ".." segments.
func Getwd() (string, error) {
	var path string
	var target string
	parent := "."
outer:
	for {
		target = parent
		// println("statting target:", target)
		dirInfo, err := os.Stat(target)
		if err != nil {
			return path, err
		}
		parent = filepath.Join("..", target)

		dir, err := os.Open(parent)
		if err != nil {
			return filepath.Clean(filepath.Join(target, path)), nil
		}
		names, err := dir.Readdirnames(-1)
		dir.Close()
		if err != nil {
			return filepath.Clean(filepath.Join(target, path)), nil
		}
		// println("NAMES: ", strings.Join(names, " "))
		for _, name := range names {
			info, err := os.Stat(filepath.Join(parent, name))
			if err != nil {
				continue
			}
			if os.SameFile(dirInfo, info) {
				path = filepath.Join(name, path)
				continue outer
			}
		}
		break
	}
	return filepath.Clean(filepath.Join("/", path)), nil
}
