package bindgen

import (
	"testing"

	"github.com/ydnar/wasm-tools-go/testdata"
	"github.com/ydnar/wasm-tools-go/wit"
)

func loadTestdata(f func(path string, res *wit.Resolve) error) error {
	return testdata.Walk(func(path string) error {
		res, err := wit.LoadJSON(path)
		if err != nil {
			return err
		}
		return f(path, res)
	}, "*.wit.json")
}

func TestGenerateTestdata(t *testing.T) {
	err := loadTestdata(func(path string, res *wit.Resolve) error {
		t.Run(testdata.Relative(path), func(t *testing.T) {
			validateGeneratedGo(t, res)
		})
		return nil
	})
	if err != nil {
		t.Error(err)
	}
}
