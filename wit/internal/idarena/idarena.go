package idarena

type Arena[T any] struct {
	ArenaID int `json:"arena_id"`
	Items   []T `json:"items"`
}

type ArenaID struct {
	ArenaID int `json:"arena_id"`
	Index   int `json:"idx"`
}
