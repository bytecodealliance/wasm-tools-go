package gen

import "testing"

func TestFileHasContent(t *testing.T) {
	positives := []File{
		{Name: "comment.go", Content: []byte("// Comment\n")},
		{Name: "package_docs.go", PackageDocs: "package documentation"},
		{Name: "header.go", Header: "// Header\n"},
		{Name: "trailer.go", Trailer: "// Trailer\n"},
		{Name: "blank_imports.go", Imports: map[string]string{"unsafe": "_"}},
		{Name: "assembly.s", Content: []byte("// Comment\n")},
	}
	for _, f := range positives {
		t.Run(f.Name, func(t *testing.T) {
			got, want := f.HasContent(), true
			if got != want {
				t.Errorf("f.HasContent(): %t, expected %t", got, want)
			}
		})
	}

	negatives := []File{
		{Name: "empty.go", GeneratedBy: "package testing"},
		{Name: "build_tag_only.go", GoBuild: "!wasip1"},
		{Name: "named_imports.go", Imports: map[string]string{"unsafe": "unsafe"}},
		{Name: "assembly.s", Content: nil},
	}
	for _, f := range negatives {
		t.Run(f.Name, func(t *testing.T) {
			got, want := f.HasContent(), false
			if got != want {
				t.Errorf("f.HasContent(): %t, expected %t", got, want)
			}
		})
	}
}

func TestFileBytes(t *testing.T) {
	pkg := NewPackage("wasm/wasi/clocks/wallclock")
	f := pkg.File("wallclock.wit.go")
	if !f.IsGo() {
		t.Errorf("file %s is not Go", f.Name)
	}
	f.Import("encoding/json")
	f.Import("io")
	_, err := f.Bytes()
	if err != nil {
		t.Error(err)
	}
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
