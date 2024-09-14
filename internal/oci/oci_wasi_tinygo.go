//go:build wasip1 || wasip2 || tinygo

package oci

import (
	"context"
	"fmt"
)

// IsOCIPath checks if a given path is an OCI path
func IsOCIPath(path string) bool {
	return false
}

func PullWIT(ctx context.Context, path, out string) error {
	return fmt.Errorf("OCI not supported on WASI")
}
