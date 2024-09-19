//go:build wasip1 || wasip2 || tinygo

package oci

import (
	"bytes"
	"context"
	"errors"
)

// IsOCIPath checks if a given path is an OCI path
func IsOCIPath(path string) bool {
	return false
}

func PullWIT(ctx context.Context, path string) (*bytes.Buffer, error) {
	return nil, errors.New("OCI not supported on WASI or TinyGo")
}
