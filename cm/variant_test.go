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

	var u1 UnsizedVariant2[struct{}, struct{}]
	var u2 UnsizedVariant2[[0]byte, struct{}]

	tests := []struct {
		v      VariantDebug
		size   uintptr
		offset uintptr
	}{
		{&u1, 1, 0},
		{&u2, 1, 0},
		{&SizedVariant2[Shape[string], string, string]{}, sizePlusAlignOf[string](), ptrSize},
		{&SizedVariant2[Shape[string], bool, string]{}, sizePlusAlignOf[string](), ptrSize},
		{&SizedVariant2[Shape[string], string, struct{}]{}, sizePlusAlignOf[string](), ptrSize},
		{&SizedVariant2[Shape[string], struct{}, string]{}, sizePlusAlignOf[string](), ptrSize},
		{&SizedVariant2[Shape[uint64], uint64, uint64]{}, 16, alignOf[uint64]()},
		{&SizedVariant2[Shape[uint64], uint32, uint64]{}, 16, alignOf[uint64]()},
		{&SizedVariant2[Shape[uint64], uint64, uint32]{}, 16, alignOf[uint64]()},
		{&SizedVariant2[Shape[uint64], uint8, uint64]{}, 16, alignOf[uint64]()},
		{&SizedVariant2[Shape[uint64], uint64, uint8]{}, 16, alignOf[uint64]()},
		{&SizedVariant2[Shape[uint32], uint8, uint32]{}, 8, alignOf[uint32]()},
		{&SizedVariant2[Shape[uint32], uint32, uint8]{}, 8, alignOf[uint32]()},
		{&SizedVariant2[Shape[[9]byte], [9]byte, uint64]{}, 24, alignOf[uint64]()},
	}

	for _, tt := range tests {
		name := typeName(tt.v)
		t.Run(name, func(t *testing.T) {
			if got, want := tt.v.Size(), tt.size; got != want {
				t.Errorf("(%s).Size() == %v, expected %v", name, got, want)
			}
			if got, want := tt.v.ValOffset(), tt.offset; got != want {
				t.Errorf("(%s).ValOffset() == %v, expected %v", name, got, want)
			}
		})
	}
}

func zeroPtr[T any]() *T {
	var zero T
	return &zero
}
