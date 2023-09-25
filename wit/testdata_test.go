package wit

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/k0kubun/pp/v3"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/ydnar/wasm-tools-go/internal/callerfs"
)

var update = flag.Bool("update", false, "update golden files")

func compareOrWrite(t *testing.T, path, golden, data string) {
	if *update {
		err := os.WriteFile(golden, []byte(data), 0644)
		if err != nil {
			t.Error(err)
		}
		return
	}
	want, err := os.ReadFile(golden)
	if err != nil {
		t.Error(err)
		return
	}
	if string(want) != data {
		dmp := diffmatchpatch.New()
		dmp.PatchMargin = 3
		diffs := dmp.DiffMain(string(want), data, false)
		t.Errorf("value for %s did not match golden value %s:\n%v", path, golden, dmp.DiffPrettyText(diffs))
	}
}

const testdataDir = "../testdata/"

func loadTestdata(f func(path string, res *Resolve) error) error {
	return filepath.WalkDir(callerfs.Path(testdataDir), func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fs.SkipDir
		}
		if !strings.HasSuffix(path, ".wit.json") && !strings.HasSuffix(path, ".wit.md.json") {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		res, err := DecodeJSON(file)
		if err != nil {
			return err
		}
		return f(path, res)
	})
}

func TestGoldenFiles(t *testing.T) {
	p := pp.New()
	p.SetExportedOnly(true)
	p.SetColoringEnabled(false)

	err := loadTestdata(func(path string, res *Resolve) error {
		t.Run(strings.TrimPrefix(path, testdataDir), func(t *testing.T) {
			data := p.Sprint(res)
			compareOrWrite(t, path, path+".golden", data)
		})
		return nil
	})

	if err != nil {
		t.Error(err)
	}
}

func TestSizeAndAlign(t *testing.T) {
	err := loadTestdata(func(path string, res *Resolve) error {
		t.Run(strings.TrimPrefix(path, testdataDir), func(t *testing.T) {
			for i, td := range res.TypeDefs {
				name := fmt.Sprintf("TypeDefs[%d]", i)
				if td.Name != nil {
					name += "#" + *td.Name
				}
				t.Run(name, func(t *testing.T) {
					defer func() {
						err := recover()
						if err != nil {
							t.Fatalf("panic: %v", err)
						}
					}()

					got, want := td.Size(), td.Kind.Size()
					if got != want {
						t.Errorf("(*TypeDef).Size(): got %d, expected %d", got, want)
					}

					got, want = td.Align(), td.Kind.Align()
					if got != want {
						t.Errorf("(*TypeDef).Align(): got %d, expected %d", got, want)
					}
				})
			}
		})
		return nil
	})
	if err != nil {
		t.Error(err)
	}
}

// TestFunctionNameConsistency tests to see if the names in the map[string]*Function in
// each [Interface] in a [Resolve] is identical to its Name field.
func TestFunctionNameConsistency(t *testing.T) {
	err := loadTestdata(func(path string, res *Resolve) error {
		t.Run(strings.TrimPrefix(path, testdataDir), func(t *testing.T) {
			for i, face := range res.Interfaces {
				if len(face.Functions) == 0 {
					continue
				}
				name := fmt.Sprintf("Interfaces[%d]", i)
				if face.Name != nil {
					name += "#" + *face.Name
				}
				t.Run(name, func(t *testing.T) {
					for name, f := range face.Functions {
						t.Run(name, func(t *testing.T) {
							if name != f.Name {
								t.Errorf("Interface.Functions[%q] != %q", name, f.Name)
							}
						})
					}
				})
			}

			for i, w := range res.Worlds {
				if len(w.Imports) == 0 && len(w.Exports) == 0 {
					continue
				}
				name := fmt.Sprintf("Worlds[%d]#%s", i, w.Name)
				t.Run(name, func(t *testing.T) {
					// A world can rename an imported function, so disable this
					// for name, item := range w.Imports {
					// 	f, ok := item.(*Function)
					// 	if !ok {
					// 		continue
					// 	}
					// 	t.Run(fmt.Sprintf("Imports[%q]==%q", name, f.Name), func(t *testing.T) {
					// 		if name != f.Name {
					// 			t.Errorf("Imports[%q] != %q", name, f.Name)
					// 		}
					// 	})
					// }

					// TODO: can a world rename an exported function?
					for name, item := range w.Exports {
						f, ok := item.(*Function)
						if !ok {
							continue
						}
						t.Run(fmt.Sprintf("Exports[%q]==%q", name, f.Name), func(t *testing.T) {
							if name != f.Name {
								t.Errorf("Exports[%q] != %q", name, f.Name)
							}
						})
					}
				})
			}
		})
		return nil
	})
	if err != nil {
		t.Error(err)
	}
}

// TestPackageFieldIsNotNil tests to ensure the Package field of all [World] and [Interface]
// values in a [Resolve] are non-nil.
func TestPackageFieldIsNotNil(t *testing.T) {
	err := loadTestdata(func(path string, res *Resolve) error {
		t.Run(strings.TrimPrefix(path, testdataDir), func(t *testing.T) {
			for i, face := range res.Interfaces {
				name := fmt.Sprintf("Interfaces[%d]", i)
				if face.Name != nil {
					name += "#" + *face.Name
				}
				t.Run(name, func(t *testing.T) {
					if face.Package == nil {
						t.Error("Package is nil")
					}
				})
			}
			for i, w := range res.Worlds {
				name := fmt.Sprintf("Worlds[%d]#%s", i, w.Name)
				t.Run(name, func(t *testing.T) {
					if w.Package == nil {
						t.Error("Package is nil")
					}
				})
			}
		})
		return nil
	})
	if err != nil {
		t.Error(err)
	}
}

// TestTypeDefNamesNotNil verifies that all [Record], [Variant], [Enum], and [Flags]
// types have a non-nil Name.
func TestTypeDefNamesNotNil(t *testing.T) {
	err := loadTestdata(func(path string, res *Resolve) error {
		t.Run(strings.TrimPrefix(path, testdataDir), func(t *testing.T) {
			for i, td := range res.TypeDefs {
				switch td.Kind.(type) {
				case *Record, *Variant, *Enum, *Flags:
				default:
					continue
				}
				name := fmt.Sprintf("TypeDefs[%d]", i)
				if td.Name != nil {
					name += "#" + *td.Name
				}
				t.Run(name, func(t *testing.T) {
					if td.Name == nil {
						t.Error("Name is nil")
					}
				})
			}
		})
		return nil
	})
	if err != nil {
		t.Error(err)
	}
}

// TestNoExportedTypeDefs verifies that any [TypeDef] instances in a [World] are
// referenced in Imports, and not Exports.
func TestNoExportedTypeDefs(t *testing.T) {
	err := loadTestdata(func(path string, res *Resolve) error {
		t.Run(strings.TrimPrefix(path, testdataDir), func(t *testing.T) {
			for i, w := range res.Worlds {
				if len(w.Imports) == 0 && len(w.Exports) == 0 {
					continue
				}
				name := fmt.Sprintf("Worlds[%d]#%s", i, w.Name)
				t.Run(name, func(t *testing.T) {
					for name, item := range w.Exports {
						if _, ok := item.(*TypeDef); ok {
							t.Errorf("found TypeDef in World.Exports: %s", name)
						}
					}
				})
			}
		})
		return nil
	})
	if err != nil {
		t.Error(err)
	}
}
