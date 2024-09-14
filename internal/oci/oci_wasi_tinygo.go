//go:build wasip1 || wasip2 || tinygo

package oci

import (
	"context"
	"errors"
)

// IsOCIPath checks if a given path is an OCI path
func IsOCIPath(path string) bool {
	return false
}

func PullWIT(ctx context.Context, path, out string) error {
	return errors.New("OCI not supported on WASI or TinyGo")
}
