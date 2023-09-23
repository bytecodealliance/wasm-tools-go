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

func TestTypeSize(t *testing.T) {
	tests := []struct {
		name  string
		v     Type
		size  uintptr
		align uintptr
	}{
		{"bool", Bool{}, 1, 1},
		{"s8", S8{}, 1, 1},
		{"u8", U8{}, 1, 1},
		{"s16", S16{}, 2, 2},
		{"u16", U16{}, 2, 2},
		{"s32", S32{}, 4, 4},
		{"u32", U32{}, 4, 4},
		{"s64", S64{}, 8, 8},
		{"u64", U64{}, 8, 8},
		{"float32", Float32{}, 4, 4},
		{"float64", Float64{}, 8, 8},
		{"char", Char{}, 4, 4},
		{"string", String{}, 8, 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			size := tt.v.Size()
			if size != tt.size {
				t.Errorf("(Type).Size(): expected %d, got %d", tt.size, size)
			}
			align := tt.v.Align()
			if align != tt.align {
				t.Errorf("(Type).Align(): expected %d, got %d", tt.align, align)
			}

		})
	}
}
