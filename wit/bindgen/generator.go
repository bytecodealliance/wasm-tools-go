// Package bindgen generates Go source code from a fully-resolved WIT package.
// It generates one or more Go packages, with functions, types, constants, and variables,
// along with the associated code to lift and lower Go types into Canonical ABI representation.
package bindgen

import (
	"path/filepath"
	"runtime"

	"github.com/ydnar/wasm-tools-go/internal/go/gen"
	"github.com/ydnar/wasm-tools-go/wit"
)

const cmPackage = "github.com/ydnar/wasm-tools-go/cm"

type generator struct {
	opts options
	res  *wit.Resolve

	// versioned is set to true if there are multiple versions of a WIT package in res,
	// which affects the generated Go package paths.
	versioned bool

	// packages are Go packages indexed on Go package paths.
	packages map[string]*gen.Package

	// witPackages map WIT identifier paths to Go packages.
	witPackages map[string]*gen.Package

	// typeDefs map [wit.TypeDef] to their Go equivalent.
	typeDefs map[*wit.TypeDef]typeDecl

	// functions map [wit.Function] to their Go equivalent.
	functions map[*wit.Function]funcDecl

	// defined represent whether a type or function has been defined.
	defined map[any]bool
}

func newGenerator(res *wit.Resolve, opts ...Option) (*generator, error) {
	g := &generator{
		packages:    make(map[string]*gen.Package),
		witPackages: make(map[string]*gen.Package),
		typeDefs:    make(map[*wit.TypeDef]typeDecl),
		functions:   make(map[*wit.Function]funcDecl),
		defined:     make(map[any]bool),
	}
	err := g.opts.apply(opts...)
	if err != nil {
		return nil, err
	}
	if g.opts.generatedBy == "" {
		_, file, _, _ := runtime.Caller(0)
		_, g.opts.generatedBy = filepath.Split(filepath.Dir(filepath.Dir(file)))
	}
	if g.opts.cmPackage == "" {
		g.opts.cmPackage = cmPackage
	}
	g.res = res
	return g, nil
}

type typeDecl struct {
	file  *gen.File // The Go file this type belongs to
	scope gen.Scope // Scope for type-local declarations like method names
	name  string    // The unique Go name for this type
}

type funcDecl struct {
	f     function // The exported Go function
	lower function // The canon lower function (go:wasmimport)
	lift  function // The canon lift function (go:wasmexport)
}

// function represents a Go function created from a Component Model function
type function struct {
	file     *gen.File // The Go file this function belongs to
	scope    gen.Scope // Scope for function-local declarations
	name     string    // The scoped unique Go name for this function (method names are scoped to receiver type)
	receiver param     // The method receiver, if any
	params   []param   // Function param(s), with unique Go name(s)
	results  []param   // Function result(s), with unique Go name(s)
}

func (f *function) isMethod() bool {
	return f.receiver.typ != nil
}

// param represents a Go function parameter or result.
// name is a unique Go name within the function scope.
type param struct {
	name string
	typ  wit.Type
}

func goFunction(file *gen.File, f *wit.Function, goName string) function {
	scope := gen.NewScope(file)
	out := function{
		file:    file,
		scope:   scope,
		name:    goName,
		params:  goParams(scope, f.Params),
		results: goParams(scope, f.Results),
	}
	if len(out.results) == 1 && out.results[0].name == "" {
		out.results[0].name = "result"
	}
	if f.IsMethod() {
		out.receiver = out.params[0]
		out.params = out.params[1:]
	}
	return out
}

func goParams(scope gen.Scope, params []wit.Param) []param {
	out := make([]param, len(params))
	for i := range params {
		out[i].name = scope.DeclareName(GoName(params[i].Name, false))
		out[i].typ = params[i].Type
	}
	return out
}
