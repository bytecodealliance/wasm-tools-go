package wit

import (
	"encoding/json"
	"errors"
	"io"
)

func DecodeJSON(r io.Reader) (*Resolve, error) {
	dec := json.NewDecoder(r)
	dec.UseNumber()

	res := &Resolve{}
	tok, err := dec.Token()
	if err != nil {
		return nil, err
	}
	if d, ok := tok.(json.Delim); !ok || d != '{' {
		return nil, errors.New("missing {")
	}
	return res, nil
}

func (res *Resolve) UnmarshalJSON(data []byte) error {
	type arena[T any] struct {
		Items unmarshalJSONFunc `json:"items"`
	}

	var proxy struct {
		Worlds     arena[World]     `json:"worlds"`
		Interfaces arena[Interface] `json:"interfaces"`
		Types      arena[TypeDef]   `json:"types"`
		Packages   arena[Package]   `json:"packages"`
	}

	proxy.Worlds.Items = unmarshalJSONFunc(func(data []byte) error {
		return nil
	})

	return json.Unmarshal(data, &proxy)
}

type items[T any] struct {
	res   *Resolve
	items *[]*T
}

func (dec *items[T]) UnmarshalJSON() error {
	return nil
}
