package oci

import (
	"os"

	"github.com/regclient/regclient/types/ref"
)

// IsOCIPath checks if a given path is an OCI path
func IsOCIPath(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return false
	}

	_, err := ref.New(path)
	return err == nil
}
