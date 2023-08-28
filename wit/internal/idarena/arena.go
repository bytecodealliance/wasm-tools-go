package idarena

type Arena[T any] struct {
	ArenaID int `json:"arena_id"`
	Items   []T `json:"items"`
}

func (a Arena[T]) Item(i int) *T {
	return &a.Items[i]
}

type ArenaID struct {
	ArenaID int `json:"arena_id"`
	Index   int `json:"idx"`
}
