// TODO: remove this once TinyGo adds runtime.Frame.Entry and reflect.StringHeader.Len is type int
//go:build !tinygo

package bindgen

import (
	"testing"

	"github.com/ydnar/wasm-tools-go/internal/relpath"
	"github.com/ydnar/wasm-tools-go/wit"
)

const testdataPath = "../../testdata"

func loadTestdata(f func(path string, res *wit.Resolve) error) error {
	return relpath.Walk(testdataPath, func(path string) error {
		res, err := wit.LoadJSON(path)
		if err != nil {
			return err
		}
		return f(path, res)
	}, "*.wit.json")
}

func TestGenerateTestdata(t *testing.T) {
	if testing.Short() {
		// t.Skip is not available in TinyGo, requires runtime.Goexit()
		return
	}
	err := loadTestdata(func(path string, res *wit.Resolve) error {
		t.Run(path, func(t *testing.T) {
			validateGeneratedGo(t, res)
		})
		return nil
	})
	if err != nil {
		t.Error(err)
	}
}
