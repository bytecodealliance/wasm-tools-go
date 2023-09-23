package wit

import (
	"fmt"
	"strings"
	"testing"
)

func TestGoldenFilesABI(t *testing.T) {
	err := loadTestdata(func(path string, res *Resolve) error {
		t.Run(strings.TrimPrefix(path, testdataDir), func(t *testing.T) {
			for i := range res.TypeDefs {
				td := res.TypeDefs[i]
				name := fmt.Sprintf("types/%d", i)
				if td.Name != nil {
					name += "/" + *td.Name
				}
				t.Run(name, func(t *testing.T) {
					defer func() {
						err := recover()
						if err != nil {
							t.Fatalf("panic: %v", err)
						}
					}()

					got, want := td.Size(), td.Kind.Size()
					if got != want {
						t.Errorf("(*TypeDef).Size(): got %d, expected %d", got, want)
					}

					got, want = td.Align(), td.Kind.Align()
					if got != want {
						t.Errorf("(*TypeDef).Align(): got %d, expected %d", got, want)
					}
				})
			}
		})
		return nil
	})
	if err != nil {
		t.Error(err)
	}
}
