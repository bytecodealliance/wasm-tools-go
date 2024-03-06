package bindgen

import (
	"github.com/ydnar/wasm-tools-go/internal/go/gen"
	"github.com/ydnar/wasm-tools-go/wit"
)

type funcDecl struct {
	f     function // The exported Go function
	lower function // The canon lower function (go:wasmimport)
	lift  function // The canon lift function (go:wasmexport)
}

// function represents a Go function created from a Component Model function
type function struct {
	file     *gen.File // The Go file this function belongs to
	scope    gen.Scope // Scope for function-local declarations
	name     string    // The scoped unique Go name for this function (method names are scoped to recevier type)
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
