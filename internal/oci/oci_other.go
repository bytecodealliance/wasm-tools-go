//go:build !wasip1 && !wasip2 && !tinygo

package oci

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/regclient/regclient"
	"github.com/regclient/regclient/types/manifest"
	"github.com/regclient/regclient/types/ref"
)

// IsOCIPath checks if a given path is an OCI path
func IsOCIPath(path string) bool {
	_, err := ref.New(path)
	return err == nil
}

// PullWIT fetches an OCI artifact from a given OCI path
// It invokes "regclient" APIs to pull the artifact and then
// processes it with `wasm-tools`.
// the WIT files are stored in the `<out>` directory
func PullWIT(ctx context.Context, path, out string) error {
	wasmTools, err := exec.LookPath("wasm-tools")
	if err != nil {
		return err
	}

	r, err := ref.New(path)
	if err != nil {
		return fmt.Errorf("failed to parse ref: %v", err)
	}

	rc := regclient.New()
	defer rc.Close(ctx, r)

	m, err := rc.ManifestGet(ctx, r)
	if err != nil {
		return fmt.Errorf("failed to get manifest: %v", err)
	}

	mi, ok := m.(manifest.Imager)
	if !ok {
		return fmt.Errorf("manifest does not support image methods")
	}

	layers, err := mi.GetLayers()
	if err != nil {
		return fmt.Errorf("failed to get layers: %v", err)
	}
	if len(layers) == 0 {
		return fmt.Errorf("no layers found in the artifact")
	}

	layer := layers[0] // assuming the WIT file is the first layer

	if err = layer.Digest.Validate(); err != nil {
		return fmt.Errorf("layer contains invalid digest: %s: %v", string(layer.Digest), err)
	}

	rdr, err := rc.BlobGet(ctx, r, layer)
	if err != nil {
		return fmt.Errorf("failed to fetch blob: %v", err)
	}
	defer rdr.Close()

	wasmCmd := exec.Command(wasmTools, "component", "wit", "--out-dir", out)
	wasmCmd.Stdin = rdr
	wasmCmd.Stdout = os.Stdout
	wasmCmd.Stderr = os.Stderr

	if err := wasmCmd.Run(); err != nil {
		return fmt.Errorf("failed to process artifact with wasm-tools: %w", err)
	}

	return nil
}
