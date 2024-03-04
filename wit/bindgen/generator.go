// Package bindgen generates Go source code from a fully-resolved WIT package.
// It generates one or more Go packages, with functions, types, constants, and variables,
// along with the associated code to lift and lower Go types into Canonical ABI representation.
package bindgen

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"

	"github.com/ydnar/wasm-tools-go/cm"
	"github.com/ydnar/wasm-tools-go/internal/codec"
	"github.com/ydnar/wasm-tools-go/internal/go/gen"
	"github.com/ydnar/wasm-tools-go/internal/stringio"
	"github.com/ydnar/wasm-tools-go/wit"
)

const (
	GoSuffix     = ".wit.go"
	BuildDefault = "!wasip1"
	cmPackage    = "github.com/ydnar/wasm-tools-go/cm"
	emptyAsm     = `// This file exists for testing this package without WebAssembly,
// allowing empty function bodies with a //go:wasmimport directive.
// See https://pkg.go.dev/cmd/compile for more information.
`
)

// Go generates one or more Go packages from [wit.Resolve] res.
// It returns any error that occurs during code generation.
func Go(res *wit.Resolve, opts ...Option) ([]*gen.Package, error) {
	g, err := newGenerator(res, opts...)
	if err != nil {
		return nil, err
	}
	return g.generate()
}

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

	// worldPackages map [wit.World] to Go packages.
	worldPackages map[*wit.World]*gen.Package

	// interfacePackages map [wit.Interface] to Go packages.
	interfacePackages map[*wit.Interface]*gen.Package

	// typeDefPackages map [wit.TypeDef] to Go packages.
	typeDefPackages map[*wit.TypeDef]*gen.Package

	// typeDefNames map [wit.TypeDef] to a defined Go name.
	typeDefNames map[*wit.TypeDef]string

	// scopes map a [wit.TypeDEf] to a [gen.Scope]. Used for method lists.
	scopes map[*wit.TypeDef]gen.Scope

	// functions map [wit.Function] to their equivalent Go identifier.
	functions map[*wit.Function]gen.Ident

	// defined represent whether a type or function has been defined.
	defined map[any]bool
}

func newGenerator(res *wit.Resolve, opts ...Option) (*generator, error) {
	g := &generator{
		packages:          make(map[string]*gen.Package),
		witPackages:       make(map[string]*gen.Package),
		worldPackages:     make(map[*wit.World]*gen.Package),
		interfacePackages: make(map[*wit.Interface]*gen.Package),
		typeDefPackages:   make(map[*wit.TypeDef]*gen.Package),
		typeDefNames:      make(map[*wit.TypeDef]string),
		scopes:            make(map[*wit.TypeDef]gen.Scope),
		functions:         make(map[*wit.Function]gen.Ident),
		defined:           make(map[any]bool),
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

func (g *generator) generate() ([]*gen.Package, error) {
	g.detectVersionedPackages()

	err := g.declareTypeDefs()
	if err != nil {
		return nil, err
	}

	err = g.defineInterfaces()
	if err != nil {
		return nil, err
	}

	err = g.defineWorlds()
	if err != nil {
		return nil, err
	}

	var packages []*gen.Package
	for _, path := range codec.SortedKeys(g.packages) {
		packages = append(packages, g.packages[path])
	}
	return packages, nil
}

func (g *generator) detectVersionedPackages() {
	if g.opts.versioned {
		g.versioned = true
		fmt.Fprintf(os.Stderr, "Generated versions for all package(s)\n")
		return
	}
	packages := make(map[string]string)
	for _, pkg := range g.res.Packages {
		id := pkg.Name
		id.Version = nil
		path := id.String()
		if packages[path] != "" && packages[path] != pkg.Name.String() {
			g.versioned = true
		} else {
			packages[path] = pkg.Name.String()
		}
	}
	if g.versioned {
		// fmt.Fprintf(os.Stderr, "Multiple versions of package(s) detected\n")
	}
}

// declareTypeDefs declares all type definitions in res.
func (g *generator) declareTypeDefs() error {
	for _, t := range g.res.TypeDefs {
		err := g.declareTypeDef(nil, "", t)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *generator) declareTypeDef(file *gen.File, goName string, t *wit.TypeDef) error {
	if goName == "" {
		if t.Name == nil {
			return nil
		}
		goName = GoName(*t.Name, true)
	}
	if file == nil {
		file = g.fileFor(typeDefOwner(t))
	}
	g.typeDefPackages[t] = file.Package
	g.typeDefNames[t] = file.Declare(goName)
	g.scopes[t] = gen.NewScope(nil)
	// fmt.Fprintf(os.Stderr, "Type:\t%s.%s\n\t%s.%s\n", owner.String(), name, decl.Package.Path, decl.Name)
	return nil
}

func typeDefOwner(t *wit.TypeDef) wit.Ident {
	var id wit.Ident
	switch owner := t.Owner.(type) {
	case *wit.World:
		id = owner.Package.Name
		id.Extension = owner.Name
	case *wit.Interface:
		id = owner.Package.Name
		if owner.Name == nil {
			id.Extension = "unknown"
		} else {
			id.Extension = *owner.Name
		}
	}
	return id
}

func (g *generator) defineInterfaces() error {
	var interfaces []*wit.Interface
	for _, i := range g.res.Interfaces {
		if i.Name != nil {
			interfaces = append(interfaces, i)
		}
	}
	// fmt.Fprintf(os.Stderr, "Generating Go for %d named interface(s)\n", len(interfaces))
	for _, i := range interfaces {
		g.defineInterface(i, *i.Name)
	}
	return nil
}

// By default, each WIT interface and world maps to a single Go package.
// Options might override the Go package, including combining multiple
// WIT interfaces and/or worlds into a single Go package.
func (g *generator) defineWorlds() error {
	// fmt.Fprintf(os.Stderr, "Generating Go for %d world(s)\n", len(g.res.Worlds))
	for _, w := range g.res.Worlds {
		g.defineWorld(w)
	}
	return nil
}

func (g *generator) defineWorld(w *wit.World) error {
	if g.defined[w] {
		return nil
	}
	if g.worldPackages[w] != nil {
		return nil
	}
	id := w.Package.Name
	id.Extension = w.Name
	pkg := g.packageFor(id)
	g.worldPackages[w] = pkg
	file := g.fileFor(id)

	var b strings.Builder
	fmt.Fprintf(&b, "Package %s represents the %s \"%s\".\n", pkg.Name, w.WITKind(), id.String())
	if w.Docs.Contents != "" {
		b.WriteString("\n")
		b.WriteString(w.Docs.Contents)
	}
	file.PackageDocs = b.String()

	// fmt.Printf("// World: %s\n\n", id.String())
	for _, name := range codec.SortedKeys(w.Imports) {
		var err error
		switch v := w.Imports[name].(type) {
		case *wit.Interface:
			err = g.defineInterface(v, name)
		case *wit.TypeDef:
			err = g.defineTypeDef(v, name)
		case *wit.Function:
			if v.IsFreestanding() {
				err = g.defineImportedFunction(v, id)
			}
		}
		if err != nil {
			return err
		}
	}

	g.defined[w] = true

	return nil
}

func (g *generator) defineInterface(i *wit.Interface, name string) error {
	if g.defined[i] {
		return nil
	}
	if g.interfacePackages[i] != nil {
		return nil
	}
	if i.Name != nil {
		name = *i.Name
	}
	id := i.Package.Name
	id.Extension = name
	pkg := g.packageFor(id)
	g.interfacePackages[i] = pkg
	file := g.fileFor(id)

	var b strings.Builder
	fmt.Fprintf(&b, "Package %s represents the %s \"%s\".\n", pkg.Name, i.WITKind(), id.String())
	if i.Docs.Contents != "" {
		b.WriteString("\n")
		b.WriteString(i.Docs.Contents)
	}
	file.PackageDocs = b.String()

	for _, name := range codec.SortedKeys(i.TypeDefs) {
		g.defineTypeDef(i.TypeDefs[name], name)
	}

	for _, name := range codec.SortedKeys(i.Functions) {
		f := i.Functions[name]
		if f.IsFreestanding() {
			g.defineImportedFunction(f, id)
		}
	}

	g.defined[i] = true

	return nil
}

func (g *generator) defineTypeDef(t *wit.TypeDef, name string) error {
	if g.defined[t] {
		return nil
	}
	if t.Name != nil {
		name = *t.Name
	}

	pkg := g.typeDefPackages[t]
	goName := g.typeDefNames[t]
	if pkg == nil || goName == "" {
		return fmt.Errorf("TypeDef %s not declared", name)
	}
	owner := typeDefOwner(t)
	file := g.fileFor(owner)

	// If an alias, get root
	root := t.Root()
	var rootOwner = typeDefOwner(root)
	var rootName = name
	if root.Name != nil {
		rootName = *root.Name
	}

	// Define the type
	var b bytes.Buffer
	stringio.Write(&b, "// ", goName, " represents the ", root.WITKind(), " \"", rootOwner.String(), "#", rootName, "\".\n")
	b.WriteString("//\n")
	if root != t {
		// Type alias
		stringio.Write(&b, "// See [", g.typeRep(file, root), "] for more information.\n")
		stringio.Write(&b, "type ", goName, " = ", g.typeRep(file, root), "\n\n")
	} else {
		b.WriteString(formatDocComments(t.Docs.Contents, false))
		b.WriteString("//\n")
		b.WriteString(formatDocComments(t.WIT(nil, ""), true))
		stringio.Write(&b, "type ", goName, " ", g.typeDefRep(file, goName, t), "\n\n")
	}

	_, err := file.Write(b.Bytes())
	if err != nil {
		return err
	}

	g.defined[t] = true

	// Define any associated functions
	if f := t.ResourceDrop(); f != nil {
		err := g.defineImportedFunction(f, owner)
		if err != nil {
			return nil
		}
	}

	if f := t.Constructor(); f != nil {
		err := g.defineImportedFunction(f, owner)
		if err != nil {
			return nil
		}
	}

	for _, f := range t.StaticFunctions() {
		err := g.defineImportedFunction(f, owner)
		if err != nil {
			return nil
		}
	}

	for _, f := range t.Methods() {
		err := g.defineImportedFunction(f, owner)
		if err != nil {
			return nil
		}
	}

	return nil
}

func (g *generator) typeDefRep(file *gen.File, goName string, t *wit.TypeDef) string {
	return g.typeDefKindRep(file, goName, t.Kind)
}

func (g *generator) typeDefKindRep(file *gen.File, goName string, kind wit.TypeDefKind) string {
	switch kind := kind.(type) {
	case *wit.Pointer:
		return g.pointerRep(file, kind)
	case wit.Type:
		return g.typeRep(file, kind)
	case *wit.Record:
		return g.recordRep(file, kind)
	case *wit.Tuple:
		return g.tupleRep(file, kind)
	case *wit.Flags:
		return g.flagsRep(file, goName, kind)
	case *wit.Enum:
		return g.enumRep(file, goName, kind)
	case *wit.Variant:
		return g.variantRep(file, goName, kind)
	case *wit.Result:
		return g.resultRep(file, kind)
	case *wit.Option:
		return g.optionRep(file, kind)
	case *wit.List:
		return g.listRep(file, kind)
	case *wit.Resource:
		return g.resourceRep(file, kind)
	case *wit.Own:
		return g.ownRep(file, kind)
	case *wit.Borrow:
		return g.borrowRep(file, kind)
	case *wit.Future:
		return "any /* TODO: *wit.Future */"
	case *wit.Stream:
		return "any /* TODO: *wit.Stream */"
	default:
		panic(fmt.Sprintf("BUG: unknown wit.TypeDefKind %T", kind)) // should never reach here
	}
}

func (g *generator) pointerRep(file *gen.File, p *wit.Pointer) string {
	return "*" + g.typeRep(file, p.Type)
}

func (g *generator) typeRep(file *gen.File, t wit.Type) string {
	switch t := t.(type) {
	// Special-case nil for the _ in result<T, _>
	case nil:
		return "struct{}"

	case *wit.TypeDef:
		if name, ok := g.typeDefNames[t]; ok {
			return file.RelativeName(g.typeDefPackages[t], name)
		}
		// FIXME: this is only valid for built-in WIT types.
		// User-defined types must be named, so the Ident check above must have succeeded.
		// See https://component-model.bytecodealliance.org/design/wit.html#built-in-types
		// and https://component-model.bytecodealliance.org/design/wit.html#user-defined-types.
		// TODO: add wit.Type.BuiltIn() method?
		return g.typeDefRep(file, "", t)
	case wit.Primitive:
		return g.primitiveRep(file, t)
	default:
		panic(fmt.Sprintf("BUG: unknown wit.Type %T", t)) // should never reach here
	}
}

func (g *generator) primitiveRep(file *gen.File, p wit.Primitive) string {
	switch p := p.(type) {
	case wit.Bool:
		return "bool"
	case wit.S8:
		return "int8"
	case wit.U8:
		return "uint8"
	case wit.S16:
		return "int16"
	case wit.U16:
		return "uint16"
	case wit.S32:
		return "int32"
	case wit.U32:
		return "uint32"
	case wit.S64:
		return "int64"
	case wit.U64:
		return "uint64"
	case wit.Float32:
		return "float32"
	case wit.Float64:
		return "float64"
	case wit.Char:
		return "rune"
	case wit.String:
		return "string"
	default:
		panic(fmt.Sprintf("BUG: unknown wit.Primitive %T", p)) // should never reach here
	}
}

func (g *generator) recordRep(file *gen.File, r *wit.Record) string {
	var b strings.Builder
	b.WriteString("struct {")
	for i, f := range r.Fields {
		if i == 0 || i > 0 && f.Docs.Contents != "" {
			b.WriteRune('\n')
		}
		b.WriteString(formatDocComments(f.Docs.Contents, false))
		stringio.Write(&b, fieldName(f.Name), " ", g.typeRep(file, f.Type), "\n")
	}
	b.WriteRune('}')
	return b.String()
}

func fieldName(name string) string {
	if name == "" {
		return ""
	}
	if name[0] >= '0' && name[0] <= '9' {
		name = "f" + name
	}
	return GoName(name, true)
}

func (g *generator) tupleRep(file *gen.File, t *wit.Tuple) string {
	var b strings.Builder
	if typ := t.Type(); typ != nil {
		stringio.Write(&b, "[", strconv.Itoa(len(t.Types)), "]", g.typeRep(file, typ))
	} else if len(t.Types) == 0 || len(t.Types) > cm.MaxTuple {
		// Force struct representation
		return g.typeDefKindRep(file, "", t.Despecialize())
	} else {
		stringio.Write(&b, file.Import(g.opts.cmPackage), ".Tuple")
		if len(t.Types) > 2 {
			b.WriteString(strconv.Itoa(len(t.Types)))
		}
		b.WriteRune('[')
		for i, typ := range t.Types {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(g.typeRep(file, typ))
		}
		b.WriteRune(']')
	}
	return b.String()
}

func (g *generator) flagsRep(file *gen.File, goName string, flags *wit.Flags) string {
	var b strings.Builder

	// FIXME: this isn't ideal
	var typ wit.Type
	size := flags.Size()
	switch size {
	case 1:
		typ = wit.U8{}
	case 2:
		typ = wit.U16{}
	case 4:
		typ = wit.U32{}
	case 8:
		typ = wit.U64{}
	default:
		panic(fmt.Sprintf("FIXME: cannot emit a flags type with %d cases", len(flags.Flags)))
	}

	b.WriteString(g.typeRep(file, typ))
	b.WriteString("\n\n")
	b.WriteString("const (\n")
	for i, flag := range flags.Flags {
		if i > 0 && flag.Docs.Contents != "" {
			b.WriteRune('\n')
		}
		b.WriteString(formatDocComments(flag.Docs.Contents, false))
		flagName := file.Declare(goName + GoName(flag.Name, true))
		b.WriteString(flagName)
		if i == 0 {
			stringio.Write(&b, " ", goName, " = 1 << iota")
		}
		b.WriteRune('\n')
	}
	b.WriteString(")\n")
	return b.String()
}

func (g *generator) enumRep(file *gen.File, goName string, e *wit.Enum) string {
	var b strings.Builder
	disc := wit.Discriminant(len(e.Cases))
	b.WriteString(g.typeRep(file, disc))
	b.WriteString("\n\n")
	b.WriteString("const (\n")
	for i, c := range e.Cases {
		if i > 0 && c.Docs.Contents != "" {
			b.WriteRune('\n')
		}
		b.WriteString(formatDocComments(c.Docs.Contents, false))
		b.WriteString(file.Declare(goName + GoName(c.Name, true)))
		if i == 0 {
			b.WriteRune(' ')
			b.WriteString(goName)
			b.WriteString(" = iota")
		}
		b.WriteRune('\n')
	}
	b.WriteString(")\n")
	return b.String()
}

func (g *generator) variantRep(file *gen.File, goName string, v *wit.Variant) string {
	// If the variant has no associated types, represent the variant as an enum.
	if e := v.Enum(); e != nil {
		return g.enumRep(file, goName, e)
	}

	disc := wit.Discriminant(len(v.Cases))
	shape := variantShape(v)
	align := variantAlign(v)

	// Emit type
	var b strings.Builder
	cm := file.Import(g.opts.cmPackage)
	stringio.Write(&b, cm, ".Variant[", g.typeRep(file, disc), ", ", g.typeRep(file, shape), ", ", g.typeRep(file, align), "]\n\n")

	// Emit cases
	for i, c := range v.Cases {
		caseNum := strconv.Itoa(i)
		caseName := GoName(c.Name, true)
		constructorName := file.Declare(goName + caseName)
		typeRep := g.typeRep(file, c.Type)

		// Emit constructor
		stringio.Write(&b, "// ", constructorName, " returns a [", goName, "] of case \"", c.Name, "\".\n")
		b.WriteString("//\n")
		b.WriteString(formatDocComments(c.Docs.Contents, false))
		stringio.Write(&b, "func ", constructorName, "(")
		dataName := "data"
		if c.Type != nil {
			stringio.Write(&b, dataName, " ", typeRep)
		}
		stringio.Write(&b, ") ", goName, " {")
		if c.Type == nil {
			stringio.Write(&b, "var ", dataName, " ", typeRep, "\n")
		}
		stringio.Write(&b, "return ", cm, ".New[", goName, "](", caseNum, ", ", dataName, ")\n")
		b.WriteString("}\n\n")

		// Emit getter
		if c.Type == nil {
			// Case without an associated type returns bool
			stringio.Write(&b, "// ", caseName, " returns true if [", goName, "] represents the variant case \"", c.Name, "\".\n")
			stringio.Write(&b, "func (self *", goName, ") ", caseName, "() bool {\n")
			stringio.Write(&b, "return ", cm, ".Tag(self) == ", caseNum)
			b.WriteString("}\n\n")
		} else {
			// Case with associated type T returns *T
			stringio.Write(&b, "// ", caseName, " returns a non-nil *[", typeRep, "] if [", goName, "] represents the variant case \"", c.Name, "\".\n")
			stringio.Write(&b, "func (self *", goName, ") ", caseName, "() *", typeRep, " {\n")
			stringio.Write(&b, "return ", cm, ".Case[", typeRep, "](self, ", caseNum, ")")
			b.WriteString("}\n\n")
		}
	}
	return b.String()
}

func (g *generator) resultRep(file *gen.File, r *wit.Result) string {
	var b strings.Builder
	b.WriteString(file.Import(g.opts.cmPackage))
	if r.OK == nil && r.Err == nil {
		b.WriteString(".Result")
	} else if r.OK == nil || (r.Err != nil && r.Err.Size() > r.OK.Size()) {
		stringio.Write(&b, ".ErrResult[", g.typeRep(file, r.OK), ", ", g.typeRep(file, r.Err), "]")
	} else {
		stringio.Write(&b, ".OKResult[", g.typeRep(file, r.OK), ", ", g.typeRep(file, r.Err), "]")
	}
	return b.String()
}

func (g *generator) optionRep(file *gen.File, o *wit.Option) string {
	var b strings.Builder
	stringio.Write(&b, file.Import(g.opts.cmPackage), ".Option[", g.typeRep(file, o.Type), "]")
	return b.String()
}

func (g *generator) listRep(file *gen.File, l *wit.List) string {
	var b strings.Builder
	stringio.Write(&b, file.Import(g.opts.cmPackage), ".List[", g.typeRep(file, l.Type), "]")
	return b.String()
}

func (g *generator) resourceRep(file *gen.File, r *wit.Resource) string {
	var b strings.Builder
	stringio.Write(&b, file.Import(g.opts.cmPackage), ".Resource")
	return b.String()
}

func (g *generator) ownRep(file *gen.File, o *wit.Own) string {
	return g.typeRep(file, o.Type)
}

func (g *generator) borrowRep(file *gen.File, b *wit.Borrow) string {
	return g.typeRep(file, b.Type)
}

func (g *generator) defineImportedFunction(f *wit.Function, owner wit.Ident) error {
	if g.defined[f] {
		return nil
	}

	file := g.fileFor(owner)

	// Setup
	core := f.CoreFunction(false)
	coreIsMethod := f.IsMethod() && core.Params[0] == f.Params[0]

	var funcName string
	var coreName string
	switch f.Kind.(type) {
	case *wit.Freestanding:
		funcName = file.Declare(GoName(f.BaseName(), true))
		coreName = file.Declare(GoName(f.BaseName(), false))

	case *wit.Constructor:
		t := f.Type().(*wit.TypeDef)
		funcName = file.Declare("New" + g.typeDefNames[t])
		coreName = file.Declare("new" + g.typeDefNames[t])

	case *wit.Static:
		t := f.Type().(*wit.TypeDef)
		funcName = file.Declare(g.typeDefNames[t] + GoName(f.BaseName(), true))
		coreName = file.Declare(GoName(*t.Name, false) + GoName(f.BaseName(), true))

	case *wit.Method:
		t := f.Type().(*wit.TypeDef)
		if t.Package().Name.Package != owner.Package {
			return fmt.Errorf("cannot emit functions in package %s to type %s", owner.Package, t.Package().Name.String())
		}
		scope := g.scopes[t]
		funcName = scope.UniqueName(GoName(f.BaseName(), true))
		if coreIsMethod {
			coreName = scope.UniqueName(GoName(f.BaseName(), false))
		} else {
			coreName = file.Declare(GoName(*t.Name, false) + funcName)
		}
	}

	// Go function parameters and results
	funcScope := gen.NewScope(file)
	funcParams := goParams(funcScope, f.Params)
	funcResults := goParams(funcScope, f.Results)
	var receiver wit.Param
	if f.IsMethod() {
		receiver = funcParams[0]
		funcParams = funcParams[1:]
	}

	// Core function parameters and results
	coreScope := gen.NewScope(file)
	coreParams := goParams(coreScope, core.Params)
	coreResults := goParams(coreScope, core.Results)
	if coreIsMethod {
		coreParams = coreParams[1:]
	}

	// TODO: make this work with compound params/results
	callerParams := slices.Clone(coreParams)
	for i := range callerParams {
		j := i - len(funcParams)
		if j < 0 {
			callerParams[i].Name = funcParams[i].Name
		} else {
			callerParams[i].Name = funcResults[j].Name
		}
	}

	// Bridging between Go and core function
	var compoundParams wit.Param
	if len(funcParams) > 0 && derefAnonRecord(coreParams[0].Type) != nil {
		name := funcScope.UniqueName("params")
		callerParams[0].Name = name
		t := derefAnonRecord(coreParams[0].Type)
		g.declareTypeDef(file, coreName+"Params", t)
		compoundParams.Name = name
		compoundParams.Type = t
	}

	var compoundResults wit.Param
	var resultsRecord *wit.Record
	if len(funcResults) > 1 && derefAnonRecord(last(coreParams).Type) != nil {
		name := funcScope.UniqueName("results")
		last(callerParams).Name = name
		t := derefAnonRecord(last(coreParams).Type)
		g.declareTypeDef(file, coreName+"Results", t)
		compoundResults.Name = name
		compoundResults.Type = t
		resultsRecord = t.Kind.(*wit.Record)
	}

	var b bytes.Buffer

	// Emit Go function
	b.WriteString(g.functionDocs(owner, funcName, f))
	b.WriteString("//go:nosplit\n")
	b.WriteString("func ")
	if f.IsMethod() {
		stringio.Write(&b, "(", receiver.Name, " ", g.typeRep(file, receiver.Type), ") ", funcName)
	} else {
		b.WriteString(funcName)
	}
	b.WriteRune('(')

	// Emit params
	for i, p := range funcParams {
		if i > 0 {
			b.WriteString(", ")
		}
		stringio.Write(&b, p.Name, " ", g.typeRep(file, p.Type))
	}
	b.WriteString(") ")

	// Emit results
	var namedResults bool
	if len(funcResults) == 1 {
		b.WriteString(g.typeRep(file, funcResults[0].Type))
	} else if len(funcResults) > 0 {
		namedResults = true
		b.WriteRune('(')
		for i, r := range funcResults {
			if i > 0 {
				b.WriteString(", ")
			}
			stringio.Write(&b, r.Name, " ", g.typeRep(file, r.Type))
		}
		b.WriteRune(')')
	}

	// Emit function body
	b.WriteString(" {\n")
	sameResults := slices.Equal(funcResults, coreResults)
	if !namedResults && !sameResults {
		for _, r := range funcResults {
			stringio.Write(&b, "var ", r.Name, " ", g.typeRep(file, r.Type), "\n")
		}
	}

	// Emit compound types
	if compoundParams.Type != nil {
		stringio.Write(&b, compoundParams.Name, " := ", g.typeRep(file, compoundParams.Type), "{ ")
		if receiver.Name != "" {
			stringio.Write(&b, receiver.Name, ", ")
		}
		for i, p := range funcParams {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(p.Name)
		}
		b.WriteString(" }\n")
	}
	if compoundResults.Type != nil {
		stringio.Write(&b, "var ", compoundResults.Name, " ", g.typeRep(file, compoundResults.Type), "\n")
	}

	// Emit call to wasmimport function
	if sameResults && len(coreResults) > 0 {
		b.WriteString("return ")
	}
	if coreIsMethod {
		stringio.Write(&b, receiver.Name, ".")
	}
	stringio.Write(&b, coreName, "(")
	for i, p := range callerParams {
		if i > 0 {
			b.WriteString(", ")
		}
		if isPointer(p.Type) {
			b.WriteRune('&')
		}
		b.WriteString(callerParams[i].Name)
	}
	b.WriteString(")\n")
	if !sameResults {
		b.WriteString("return ")
		if resultsRecord != nil {
			for i, f := range resultsRecord.Fields {
				if i > 0 {
					b.WriteString(", ")
				}
				stringio.Write(&b, compoundResults.Name, ".", fieldName(f.Name))
			}
		} else {
			for i, r := range funcResults {
				if i > 0 {
					b.WriteString(", ")
				}
				b.WriteString(r.Name)
			}
		}
		b.WriteRune('\n')
	}
	b.WriteString("}\n\n")

	// Emit wasmimport function
	stringio.Write(&b, "//go:wasmimport ", owner.String(), " ", f.Name, "\n")
	b.WriteString("//go:noescape\n")
	b.WriteString("func ")
	if coreIsMethod {
		stringio.Write(&b, "(", receiver.Name, " ", g.typeRep(file, receiver.Type), ") ", coreName)
	} else {
		b.WriteString(coreName)
	}
	b.WriteRune('(')

	// Emit params
	for i, p := range coreParams {
		if i > 0 {
			b.WriteString(", ")
		}
		stringio.Write(&b, p.Name, " ", g.typeRep(file, p.Type))
	}
	b.WriteString(") ")

	// Emit results
	if len(coreResults) == 1 {
		b.WriteString(g.typeRep(file, coreResults[0].Type))
	} else if len(coreResults) > 0 {
		b.WriteRune('(')
		for i, r := range coreResults {
			if i > 0 {
				b.WriteString(", ")
			}
			stringio.Write(&b, r.Name, " ", g.typeRep(file, r.Type))
		}
		b.WriteRune(')')
	}
	b.WriteString("\n\n")

	// Emit shared types
	if t, ok := compoundParams.Type.(*wit.TypeDef); ok {
		goName := g.typeDefNames[t]
		stringio.Write(&b, "// ", goName, " represents the flattened function params for [", coreName, "].\n")
		stringio.Write(&b, "// See the Canonical ABI flattening rules for more information.\n")
		stringio.Write(&b, "type ", goName, " ", g.typeDefRep(file, goName, t), "\n\n")
	}

	if t, ok := compoundResults.Type.(*wit.TypeDef); ok {
		goName := g.typeDefNames[t]
		stringio.Write(&b, "// ", goName, " represents the flattened function results for [", coreName, "].\n")
		stringio.Write(&b, "// See the Canonical ABI flattening rules for more information.\n")
		stringio.Write(&b, "type ", goName, " ", g.typeDefRep(file, goName, t), "\n\n")
	}

	// Write to file
	file.Write(b.Bytes())

	g.defined[f] = true

	return g.ensureEmptyAsm(file.Package)
}

func last[S ~[]E, E any](s S) *E {
	if len(s) == 0 {
		return nil
	}
	return &s[len(s)-1]
}

func isPointer(t wit.Type) bool {
	if td, ok := t.(*wit.TypeDef); ok {
		if _, ok := td.Kind.(*wit.Pointer); ok {
			return true
		}
	}
	return false
}

func derefTypeDef(t wit.Type) *wit.TypeDef {
	if td, ok := t.(*wit.TypeDef); ok {
		if p, ok := td.Kind.(*wit.Pointer); ok {
			if td, ok := p.Type.(*wit.TypeDef); ok {
				return td
			}
		}
	}
	return nil
}

func derefAnonRecord(t wit.Type) *wit.TypeDef {
	if td := derefTypeDef(t); td != nil && td.Name == nil && td.Owner == nil {
		if _, ok := td.Kind.(*wit.Record); ok {
			return td
		}
	}
	return nil
}

// goParams adapts WIT params to Go params, with a special case for the unnamed single result.
// It accepts a scope and string map to map WIT names to Go names.
// The resulting slice of [wit.Param] replaces the WIT names with valid, scoped Go names.
func goParams(scope gen.Scope, params []wit.Param) []wit.Param {
	params = slices.Clone(params)
	if len(params) == 1 && params[0].Name == "" {
		params[0].Name = "result"
	}
	for i := range params {
		p := &params[i]
		p.Name = scope.UniqueName(GoName(p.Name, false))
	}
	return params
}

func (g *generator) functionDocs(owner wit.Ident, goName string, f *wit.Function) string {
	var b strings.Builder
	kind := f.WITKind()
	if f.IsAdmin() {
		kind = "the Canonical ABI function"
	}
	stringio.Write(&b, "// ", goName, " represents ", kind, " \"")
	if f.IsFreestanding() {
		stringio.Write(&b, owner.String(), "#", f.Name)
	} else {
		stringio.Write(&b, f.BaseName())
	}
	b.WriteString("\".\n")
	if f.Docs.Contents != "" {
		b.WriteString("//\n")
		b.WriteString(formatDocComments(f.Docs.Contents, false))
	}
	b.WriteString("//\n")
	if !f.IsAdmin() {
		w := strings.TrimSuffix(f.WIT(nil, f.BaseName()), ";")
		b.WriteString(formatDocComments(w, true))
	}
	return b.String()
}

func (g *generator) ensureEmptyAsm(pkg *gen.Package) error {
	f := pkg.File("empty.s")
	if len(f.Content) > 0 {
		return nil
	}
	_, err := f.Write([]byte(emptyAsm))
	return err
}

func (g *generator) fileFor(id wit.Ident) *gen.File {
	pkg := g.packageFor(id)
	file := pkg.File(id.Extension + GoSuffix)
	file.GeneratedBy = g.opts.generatedBy
	file.Build = BuildDefault
	return file
}

func (g *generator) packageFor(id wit.Ident) *gen.Package {
	// Find existing
	pkg := g.witPackages[id.String()]
	if pkg != nil {
		return pkg
	}

	// Create a new package
	path := id.Namespace + "/" + id.Package + "/" + id.Extension
	if g.opts.packageRoot != "" && g.opts.packageRoot != "std" {
		path = g.opts.packageRoot + "/" + path
	}
	name := id.Extension
	if g.versioned && id.Version != nil {
		path += "/v" + id.Version.String()
	}

	// TODO: write tests for this
	name = GoPackageName(name)
	// Ensure local name doesn’t conflict with Go keywords or predeclared identifiers
	if gen.UniqueName(name, gen.IsReserved) != name {
		// Try with package prefix, like error -> ioerror
		name = FlatName(id.Package + name)
		if gen.UniqueName(name, gen.IsReserved) != name {
			// Try with namespace prefix, like ioerror -> wasiioerror
			name = gen.UniqueName(FlatName(id.Namespace+name), gen.IsReserved)
		}
	}

	pkg = gen.NewPackage(path + "#" + name)
	g.packages[pkg.Path] = pkg
	g.witPackages[id.String()] = pkg

	return pkg
}
