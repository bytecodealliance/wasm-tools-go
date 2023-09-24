package wit

import (
	"fmt"
	"math"
	"strings"
	"testing"
)

func TestAlign(t *testing.T) {
	tests := []struct {
		ptr   uintptr
		align uintptr
		want  uintptr
	}{
		{0, 1, 0}, {0, 2, 0}, {0, 4, 0}, {0, 8, 0},
		{1, 1, 1}, {1, 2, 2}, {1, 4, 4}, {1, 8, 8},
		{2, 1, 2}, {2, 2, 2}, {2, 4, 4}, {2, 8, 8},
		{3, 1, 3}, {3, 2, 4}, {3, 4, 4}, {3, 8, 8},
		{4, 1, 4}, {4, 2, 4}, {4, 4, 4}, {4, 8, 8},
		{5, 1, 5}, {5, 2, 6}, {5, 4, 8}, {5, 8, 8},
		{6, 1, 6}, {6, 2, 6}, {6, 4, 8}, {6, 8, 8},
		{7, 1, 7}, {7, 2, 8}, {7, 4, 8}, {7, 8, 8},
		{8, 1, 8}, {8, 2, 8}, {8, 4, 8}, {8, 8, 8},
		{9, 1, 9}, {9, 2, 10}, {9, 4, 12}, {9, 8, 16},
		{10, 1, 10}, {10, 2, 10}, {10, 4, 12}, {10, 8, 16},
	}
	for _, tt := range tests {
		name := fmt.Sprintf("%d,%d=%d", tt.ptr, tt.align, tt.want)
		t.Run(name, func(t *testing.T) {
			got := Align(tt.ptr, tt.align)
			if got != tt.want {
				t.Errorf("Align(%d, %d): expected %d, got %d", tt.ptr, tt.align, tt.want, got)
			}
		})
	}
}

func TestDiscriminant(t *testing.T) {
	tests := []struct {
		n    int
		want Type
	}{
		{0, U8{}}, {1, U8{}}, {5, U8{}}, {8, U8{}}, {255, U8{}}, {256, U8{}},
		{257, U16{}}, {10000, U16{}}, {32768, U16{}}, {65536, U16{}},
		{65537, U32{}}, {1 << 24, U32{}}, {math.MaxInt32, U32{}}, {math.MaxUint32, U32{}},
	}
	for _, tt := range tests {
		name := fmt.Sprintf("%d", tt.n)
		t.Run(name, func(t *testing.T) {
			got := Discriminant(tt.n)
			if got != tt.want {
				t.Errorf("Discriminant(%d): expected %T, got %T", tt.n, tt.want, got)
			}
		})
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

func TestSizeAndAlign(t *testing.T) {
	err := loadTestdata(func(path string, res *Resolve) error {
		t.Run(strings.TrimPrefix(path, testdataDir), func(t *testing.T) {
			for i, td := range res.TypeDefs {
				name := fmt.Sprintf("TypeDefs[%d]", i)
				if td.Name != nil {
					name += "#" + *td.Name
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

// TestFunctionNameConsistency tests to see if the names in the map[string]*Function in
// each [Interface] in a [Resolve] is identical to its Name field.
func TestFunctionNameConsistency(t *testing.T) {
	err := loadTestdata(func(path string, res *Resolve) error {
		t.Run(strings.TrimPrefix(path, testdataDir), func(t *testing.T) {
			for i, face := range res.Interfaces {
				if len(face.Functions) == 0 {
					continue
				}
				name := fmt.Sprintf("Interfaces[%d]", i)
				if face.Name != nil {
					name += "#" + *face.Name
				}
				t.Run(name, func(t *testing.T) {
					for name, f := range face.Functions {
						t.Run(name, func(t *testing.T) {
							if name != f.Name {
								t.Errorf("Interface.Functions[%q] != %q", name, f.Name)
							}
						})
					}
				})
			}

			for i, w := range res.Worlds {
				if len(w.Imports) == 0 && len(w.Exports) == 0 {
					continue
				}
				name := fmt.Sprintf("Worlds[%d]#%s", i, w.Name)
				t.Run(name, func(t *testing.T) {
					// A world can rename an imported function, so disable this
					// for name, item := range w.Imports {
					// 	f, ok := item.(*Function)
					// 	if !ok {
					// 		continue
					// 	}
					// 	t.Run(fmt.Sprintf("Imports[%q]==%q", name, f.Name), func(t *testing.T) {
					// 		if name != f.Name {
					// 			t.Errorf("Imports[%q] != %q", name, f.Name)
					// 		}
					// 	})
					// }

					// TODO: can a world rename an exported function?
					for name, item := range w.Exports {
						f, ok := item.(*Function)
						if !ok {
							continue
						}
						t.Run(fmt.Sprintf("Exports[%q]==%q", name, f.Name), func(t *testing.T) {
							if name != f.Name {
								t.Errorf("Exports[%q] != %q", name, f.Name)
							}
						})
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
