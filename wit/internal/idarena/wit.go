package idarena

type WIT struct {
	Worlds     Arena[World]
	Interfaces Arena[Interface]
	Packages   Arena[Package]
}

type Package struct {
	Docs Docs `json:"docs"`
}

type World struct {
	Docs Docs `json:"docs"`
}

type Interface struct {
	Name string `json:"name"`
	Docs Docs   `json:"docs"`
}
type Docs struct {
	Contents *string `json:"contents"`
}

type Type struct {
	Name string `json:"name"`
	Docs Docs   `json:"docs"`
	// TODO: Owner // typed ArenaID
}
