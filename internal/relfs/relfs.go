package relfs

import (
	"os"
	"path/filepath"
)

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
