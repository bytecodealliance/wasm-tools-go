package bindgen

import (
	"go/token"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ydnar/wasm-tools-go/internal/callerfs"
	"github.com/ydnar/wasm-tools-go/internal/codec"
	"github.com/ydnar/wasm-tools-go/internal/go/gen"
	"github.com/ydnar/wasm-tools-go/wit"
	"golang.org/x/tools/go/packages"
)

func TestGenerator(t *testing.T) {
	res, err := wit.LoadJSON(callerfs.Path("../../testdata/wasi/cli.wit.json"))
	if err != nil {
		t.Error(err)
		return
	}
	validateGeneratedGo(t, res)
}

// validateGeneratedGo loads the Go package(s) generated
func validateGeneratedGo(t *testing.T, res *wit.Resolve) {
	out := callerfs.Path(".")
	pkgPath, err := gen.PackagePath(out)
	if err != nil {
		t.Error(err)
		return
	}

	pkgs, err := Go(res,
		GeneratedBy("test"),
		PackageRoot(pkgPath),
	)
	if err != nil {
		t.Error(err)
		return
	}

	pkgMap := make(map[string]*gen.Package)

	cfg := &packages.Config{
		Mode:    packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles | packages.NeedImports | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedTypesSizes,
		Fset:    token.NewFileSet(),
		Overlay: make(map[string][]byte),
	}

	for _, pkg := range pkgs {
		pkgMap[pkg.Path] = pkg
		dir := filepath.Join(out, strings.TrimPrefix(pkg.Path, pkgPath))
		for _, file := range pkg.Files {
			path := filepath.Join(dir, file.Name)
			src, err := file.Bytes()
			if err != nil {
				t.Error(err)
			}
			cfg.Overlay[path] = src // Keep unformatted file for more testing
		}
	}

	goPackages, err := packages.Load(cfg, codec.Keys(pkgMap)...)
	if err != nil {
		t.Error(err)
		return
	}

	for _, goPkg := range goPackages {
		pkg := pkgMap[goPkg.PkgPath]
		if pkg == nil {
			t.Logf("Skipped package: %s", goPkg.PkgPath)
			continue
		}

		// Verify number of files
		count := len(goPkg.OtherFiles)
		// t.Logf("Go package: %s %t", goPkg.PkgPath, goPkg.Types.Complete())
		for _, f := range goPkg.GoFiles {
			count++
			base := filepath.Base(f)
			// t.Logf("Go file: %s", base)
			if pkg.Files[base] == nil {
				t.Errorf("unknown file in package %s: %s", pkg.Path, base)
			}
		}
		if count != len(pkg.Files) {
			t.Errorf("%d files in package %s; expected %d", count, pkg.Path, len(pkg.Files))
		}

		// Verify generated names
		for id, def := range goPkg.TypesInfo.Defs {
			if def == nil || def.Parent() != goPkg.Types.Scope() {
				continue
			}
			// t.Logf("Def: %s", id.String())
			if !pkg.HasName(id.String()) {
				t.Errorf("name %s not found in generated package %s", id.String(), pkg.Path)
			}
		}
	}
}
