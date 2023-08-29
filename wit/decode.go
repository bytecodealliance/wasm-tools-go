package wit

import (
	"encoding/json"
	"io"

	"github.com/ydnar/wit-bindgen-go/internal/wjson"
)

func DecodeJSON(r io.Reader) (*Resolve, error) {
	res := &Resolve{}

	dec := json.NewDecoder(r)
	dec.UseNumber()

	err := decodeResolve(dec, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func decodeResolve(dec *json.Decoder, res *Resolve) error {
	return wjson.DecodeObject(dec, func(worldKey string) error {
		switch worldKey {
		case "worlds":
			return decodeResolveItem(dec, res, &res.Worlds)
		case "interfaces":
			return decodeResolveItem(dec, res, &res.Interfaces)
		case "types":
			return decodeResolveItem(dec, res, &res.Types)
		case "packages":
			return decodeResolveItem(dec, res, &res.Packages)
		default:
			return wjson.Ignore(dec)
		}
	})
}

func decodeResolveItem[S ~[]*E, E any](dec *json.Decoder, res *Resolve, s *S) error {
	return wjson.DecodeObject(dec, func(key string) error {
		switch key {
		case "items":
			return wjson.DecodeArray(dec, func(i int) error {
				// TODO: decode an item
				e := element(s, i)
				_ = e
				return nil
			})
		default:
			return wjson.Ignore(dec)
		}
	})
}
