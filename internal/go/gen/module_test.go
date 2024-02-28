package gen

import (
	"os"
	"testing"

	"github.com/ydnar/wasm-tools-go/internal/callerfs"
)

func TestPackagePath(t *testing.T) {
	root := callerfs.Path(".")
	got, err := PackagePath(root)
	if err != nil {
		t.Error(err)
	}
	want := "github.com/ydnar/wasm-tools-go/internal/go/gen"
	if got != want {
		t.Errorf("PackagePath(%q): got %s, expected %s", root, got, want)
	}

	tmp := os.TempDir()
	_, err = PackagePath(tmp)
	if err == nil {
		t.Errorf("PackagePath(%q): expected error, got nil", tmp)
	}
}
