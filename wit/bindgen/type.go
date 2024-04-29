package bindgen

import (
	"github.com/ydnar/wasm-tools-go/internal/go/gen"
)

type typeDecl struct {
	file  *gen.File // The Go file this type belongs to
	scope gen.Scope // Scope for type-local declarations like method names
	name  string    // The unique Go name for this type
}
