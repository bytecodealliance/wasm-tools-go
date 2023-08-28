package idarena

type Resolve struct {
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

type WorldItem struct {
	Interface *ArenaID  `json:"Interface,omitempty"`
	Function  *Function `json:"Function,omitempty"`
	Type      *ArenaID  `json:"Type,omitempty"`
}

type Interface struct {
	Docs      Docs                `json:"docs"`
	Name      *string             `json:"name"`
	Types     map[string]ArenaID  `json:"types"`
	Functions map[string]Function `json:"functions"`
	Package   *ArenaID            `json:"package"`
}

/*
pub enum Type {
    Bool,
    U8,
    U16,
    U32,
    U64,
    S8,
    S16,
    S32,
    S64,
    Float32,
    Float64,
    Char,
    String,
    Id(TypeId),
}
*/

type Type interface {
	isType()
}

type BoolType struct{}
type U8Type struct{}

type Function struct {
	Docs    Docs
	Name    string
	Kind    FunctionKind
	Params  []Param // Vec<(String, Type)>;
	Results Results // enum
}

type FunctionKind struct {
	Freestanding bool     `json:"freestanding,omitempty"`
	Method       *ArenaID `json:"method,omitempty"`
	Static       *ArenaID `json:"static,omitempty"`
	Constructor  *ArenaID `json:"constructor,omitempty"`
}

type Docs struct {
	Contents *string `json:"contents"`
}

type Param struct {
	Name string
	Type Type
}

type Results struct {
	Named []Param
	Anon  *ArenaID
}

func foo() {

}
