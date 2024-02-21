package gen

// UniqueName tests name against filters and modifies name until all filters return false.
// Use IsReserved to filter out Go keywords and predeclared identifiers.
//
// Exported names that start with a capital letter will be appended with an underscore.
// Non-exported names are prefixed with an underscore.
func UniqueName(name string, filters ...func(string) bool) string {
	filter := func(name string) bool {
		for _, f := range filters {
			if f(name) {
				return true
			}
		}
		return false
	}
	for filter(name) {
		name += "_"
	}
	return name
}

// HasKey returns a function for map m that tests presence of key k.
// The function returns true if m contains k, and false otherwise.
func HasKey[M ~map[K]V, K comparable, V any](m M) func(k K) bool {
	return func(k K) bool {
		_, ok := m[k]
		return ok
	}
}

// Scope represents a Go name scope, like a package, file, interface, struct, or function blocks.
type Scope interface {
	// HasName returns true if this scope or any of its parent scopes contains name.
	HasName(name string) bool

	// UniqueName modifies name if necessary and declares it within this scope.
	// It returns the unique generated name.
	// Subsequent calls to HasName with the returned name will return true.
	// Subsequent calls to UniqueName will return a different name.
	UniqueName(name string) string
}

type scope struct {
	parent Scope
	names  map[string]bool
}

// NewScope returns an initialized [Scope] that's ready to use.
// If parent is nil, [Reserved] will be used.
func NewScope(parent Scope) Scope {
	if parent == nil {
		parent = Reserved()
	}
	return &scope{
		parent: parent,
		names:  make(map[string]bool),
	}
}

func (s *scope) HasName(name string) bool {
	return s.names[name] || s.parent.HasName(name)
}

func (s *scope) UniqueName(name string) string {
	name = UniqueName(name, s.HasName)
	s.names[name] = true
	return name
}

type reservedScope struct{}

// Reserved returns a preset [Scope] with the default [Go keywords] and [predeclared identifiers].
// Calls to its UniqueName method will panic, as the scope is immutable.
//
// [Go keywords]: https://go.dev/ref/spec#Keywords
// [predeclared identifiers]: https://go.dev/ref/spec#Predeclared_identifiers
func Reserved() Scope {
	return reservedScope{}
}

func (reservedScope) HasName(name string) bool {
	return IsReserved(name)
}

func (reservedScope) UniqueName(name string) string {
	panic("cannot add a name to reserved scope")
}

// IsReserved returns true for any name that is a Go keyword or predeclared identifier.
func IsReserved(name string) bool {
	return reserved[name]
}

var reserved = mapWords(
	// Keywords
	"break",
	"case",
	"chan",
	"const",
	"continue",
	"default",
	"defer",
	"else",
	"fallthrough",
	"for",
	"func",
	"go",
	"goto",
	"if",
	"import",
	"interface",
	"map",
	"package",
	"range",
	"return",
	"select",
	"struct",
	"switch",
	"type",
	"var",

	// Types
	"any",
	"bool",
	"byte",
	"comparable",
	"complex64",
	"complex128",
	"error",
	"float32",
	"float64",
	"int",
	"int8",
	"int16",
	"int32",
	"int64",
	"rune",
	"string",
	"uint",
	"uint8",
	"uint16",
	"uint32",
	"uint64",
	"uintptr",

	// Constants
	"true",
	"false",
	"iota",

	// Zero value
	"nil",

	// Functions
	"append",
	"cap",
	"clear",
	"close",
	"complex",
	"copy",
	"delete",
	"imag",
	"len",
	"make",
	"max",
	"min",
	"new",
	"panic",
	"print",
	"println",
	"real",
	"recover",
)

// Initialisms is a set of common initialisms.
var Initialisms = mapWords(
	"acl",
	"abi",
	"api",
	"ascii",
	"cabi",
	"cpu",
	"css",
	"cwd",
	"dns",
	"eof",
	"fifo",
	"guid",
	"html",
	"http",
	"https",
	"id",
	"imap",
	"io",
	"ip",
	"js",
	"json",
	"lhs",
	"mime",
	"posix",
	"qps",
	"ram",
	"rhs",
	"rpc",
	"sla",
	"smtp",
	"sql",
	"ssh",
	"tcp",
	"tls",
	"ttl",
	"tty",
	"udp",
	"ui",
	"uid",
	"uuid",
	"uri",
	"url",
	"utf8",
	"vm",
	"xml",
	"xmpp",
	"xsrf",
	"xss",
)

func mapWords(words ...string) map[string]bool {
	m := make(map[string]bool, len(words))
	for _, word := range words {
		m[word] = true
	}
	return m
}
