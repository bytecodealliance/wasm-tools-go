package idarena

type WIT struct {
	Worlds     Arena[World]     `json:"worlds"`
	Interfaces Arena[Interface] `json:"interfaces"`
	Packages   Arena[Package]   `json:"packages"`
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

type Type struct {
	Name string `json:"name"`
	Docs Docs   `json:"docs"`
	// TODO: Owner // typed ArenaID
}

type Docs struct {
	Contents *string `json:"contents"`
}
