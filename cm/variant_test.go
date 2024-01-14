package cm

import (
	"testing"
	"unsafe"
)

var (
	// TODO: nuke Variant2 interface?
	_ Variant2[struct{}, struct{}] = zeroPtr[UnsizedVariant2[struct{}, struct{}]]()
	_ Variant2[struct{}, struct{}] = zeroPtr[UnsizedVariant2[struct{}, struct{}]]()
	_ Variant2[string, bool]       = &SizedVariant2[Shape[string], string, bool]{}
	_ Variant2[bool, string]       = &SizedVariant2[Shape[string], bool, string]{}
)

func TestVariantLayout(t *testing.T) {
	// 8 on 64-bit, 4 on 32-bit
	ptrSize := unsafe.Sizeof(uintptr(0))

	tests := []struct {
		name   string
		v      VariantDebug
		size   uintptr
		offset uintptr
	}{
		{"variant { _; _ }", UnsizedVariant2[struct{}, struct{}](false), 1, 0},
		{"variant { [0]u8; _ }", UnsizedVariant2[[0]byte, struct{}](false), 1, 0},
		{"variant { string; string }", SizedVariant2[Shape[string], string, string]{}, sizePlusAlignOf[string](), ptrSize},
		{"variant { bool; string }", SizedVariant2[Shape[string], bool, string]{}, sizePlusAlignOf[string](), ptrSize},
		{"variant { string; _ }", SizedVariant2[Shape[string], string, struct{}]{}, sizePlusAlignOf[string](), ptrSize},
		{"variant { _; _ }", SizedVariant2[Shape[string], struct{}, string]{}, sizePlusAlignOf[string](), ptrSize},
		{"variant { u64; u64 }", SizedVariant2[Shape[uint64], uint64, uint64]{}, 16, alignOf[uint64]()},
		{"variant { u32; u64 }", SizedVariant2[Shape[uint64], uint32, uint64]{}, 16, alignOf[uint64]()},
		{"variant { u64; u32 }", SizedVariant2[Shape[uint64], uint64, uint32]{}, 16, alignOf[uint64]()},
		{"variant { u8; u64 }", SizedVariant2[Shape[uint64], uint8, uint64]{}, 16, alignOf[uint64]()},
		{"variant { u64; u8 }", SizedVariant2[Shape[uint64], uint64, uint8]{}, 16, alignOf[uint64]()},
		{"variant { u8; u32 }", SizedVariant2[Shape[uint32], uint8, uint32]{}, 8, alignOf[uint32]()},
		{"variant { u32; u8 }", SizedVariant2[Shape[uint32], uint32, uint8]{}, 8, alignOf[uint32]()},
		{"variant { [9]u8, u64 }", SizedVariant2[Shape[[9]byte], [9]byte, uint64]{}, 24, alignOf[uint64]()},
	}

	for _, tt := range tests {
		typ := typeName(tt.v)
		t.Run(tt.name, func(t *testing.T) {
			if got, want := tt.v.Size(), tt.size; got != want {
				t.Errorf("(%s).Size() == %v, expected %v", typ, got, want)
			}
			if got, want := tt.v.ValOffset(), tt.offset; got != want {
				t.Errorf("(%s).ValOffset() == %v, expected %v", typ, got, want)
			}
		})
	}
}
