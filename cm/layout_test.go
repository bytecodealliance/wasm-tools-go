package cm

import (
	"testing"
	"unsafe"
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

func TestResultLayout(t *testing.T) {
	// 8 on 64-bit, 4 on 32-bit
	ptrSize := unsafe.Sizeof(uintptr(0))

	tests := []struct {
		r      ResultDebug
		size   uintptr
		offset uintptr
	}{
		{&UntypedResult{}, 1, 0},

		{&UnsizedResult[struct{}, struct{}]{}, 1, 0},
		{&UnsizedResult[[0]byte, struct{}]{}, 1, 0},

		{&SizedResult[Shape[string], string, string]{}, sizePlusAlignOf[string](), ptrSize},
		{&SizedResult[Shape[string], bool, string]{}, sizePlusAlignOf[string](), ptrSize},
		{&SizedResult[Shape[string], string, struct{}]{}, sizePlusAlignOf[string](), ptrSize},
		{&SizedResult[Shape[string], struct{}, string]{}, sizePlusAlignOf[string](), ptrSize},
		{&SizedResult[Shape[uint64], uint64, uint64]{}, 16, alignOf[uint64]()},
		{&SizedResult[Shape[uint64], uint32, uint64]{}, 16, alignOf[uint64]()},
		{&SizedResult[Shape[uint64], uint64, uint32]{}, 16, alignOf[uint64]()},
		{&SizedResult[Shape[uint64], uint8, uint64]{}, 16, alignOf[uint64]()},
		{&SizedResult[Shape[uint64], uint64, uint8]{}, 16, alignOf[uint64]()},
		{&SizedResult[Shape[uint32], uint8, uint32]{}, 8, alignOf[uint32]()},
		{&SizedResult[Shape[uint32], uint32, uint8]{}, 8, alignOf[uint32]()},
		{&SizedResult[Shape[[9]byte], [9]byte, uint64]{}, 24, alignOf[uint64]()},

		{&OKSizedResult[string, struct{}]{}, sizePlusAlignOf[string](), ptrSize},
		{&OKSizedResult[string, bool]{}, sizePlusAlignOf[string](), ptrSize},
		{&OKSizedResult[[9]byte, uint64]{}, sizePlusAlignOf[string](), alignOf[uint64]()},

		{&ErrSizedResult[struct{}, string]{}, sizePlusAlignOf[string](), ptrSize},
		{&ErrSizedResult[bool, string]{}, sizePlusAlignOf[string](), ptrSize},
		{&ErrSizedResult[uint64, [9]byte]{}, sizePlusAlignOf[string](), alignOf[uint64]()},
	}

	for _, tt := range tests {
		name := typeName(tt.r)
		t.Run(name, func(t *testing.T) {
			if got, want := tt.r.Size(), tt.size; got != want {
				t.Errorf("(%s).Size() == %v, expected %v", name, got, want)
			}
			if got, want := tt.r.ValOffset(), tt.offset; got != want {
				t.Errorf("(%s).ValOffset() == %v, expected %v", name, got, want)
			}
		})
	}
}

func sizePlusAlignOf[T any]() uintptr {
	var v T
	return unsafe.Sizeof(v) + unsafe.Alignof(v)
}

func alignOf[T any]() uintptr {
	var v T
	return unsafe.Alignof(v)
}
