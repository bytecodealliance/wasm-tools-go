package types_test

import (
	"testing"

	"github.com/ydnar/wasm-tools-go/wasi/cli/filesystem/types"
)

func TestDescriptorFlags(t *testing.T) {
	var flags types.DescriptorFlags
	flags.Set(types.DescriptorFlagRead)
	flags.Set(types.DescriptorFlagWrite)
}
