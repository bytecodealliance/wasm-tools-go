package wit

type WIT struct {
	Worlds         []World
	Interfaces     []Interface
	Packages       []Package
	PackagesByName map[string]*Package
}

type World struct{}

type Interface struct{}

type Package struct{}
