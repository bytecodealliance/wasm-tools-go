package gen

import (
	"os"
	"testing"

	"github.com/bytecodealliance/wasm-tools-go/internal/relpath"
)

func TestPackagePath(t *testing.T) {
	wd, err := relpath.Getwd()
	if err != nil {
		t.Error(err)
	}

	got, err := PackagePath(wd)
	if err != nil {
		t.Error(err)
	}
	want := "github.com/bytecodealliance/wasm-tools-go/internal/go/gen"
	if got != want {
		t.Errorf("PackagePath(%q): got %s, expected %s", wd, got, want)
	}

	tmp := os.TempDir()
	_, err = PackagePath(tmp)
	if err == nil {
		t.Errorf("PackagePath(%q): expected error, got nil", tmp)
	}
}
