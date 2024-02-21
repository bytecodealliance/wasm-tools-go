package gen

import (
	"fmt"
	"testing"
)

func TestFileBytes(t *testing.T) {
	pkg := NewPackage("wasm/wasi/clocks/wallclock")
	f := pkg.File("wallclock.wit.go")
	if !f.IsGo() {
		t.Errorf("file %s is not Go", f.Name)
	}
	f.Import("encoding/json")
	f.Import("io")
	b, err := f.Bytes()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(b))
}

func TestFileAddImport(t *testing.T) {
	pkg := NewPackage("wasm/wasi/clocks/wallclock")
	f := pkg.File("wallclock.wit.go")
	if !f.IsGo() {
		t.Errorf("file %s is not Go", f.Name)
	}

	tests := []struct {
		path string
		name string
	}{
		{"encoding/json", "json"},
		{"encoding/xml", "xml"},
		{"example/error", "error_"},
		{"example/error", "error_"},
		{"example/foo#example_foo", "example_foo"},
		{"example/foo#example_foo2", "example_foo"},
		{"example/chan", "chan_"},
		{"example/chan", "chan_"},
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := f.Import(tt.path)
			if got != tt.name {
				t.Errorf("AddImport(%q): %s, expected %s", tt.path, got, tt.name)
			}
		})
	}
}
