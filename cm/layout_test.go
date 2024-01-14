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
		name   string
		r      ResultDebug
		size   uintptr
		offset uintptr
	}{
		{"result", &UntypedResult{}, 1, 0},

		{"result<_, _>", &UnsizedResult[struct{}, struct{}]{}, 1, 0},
		{"result<[0]byte, _>", &UnsizedResult[[0]byte, struct{}]{}, 1, 0},

		{"result<string, string>", &SizedResult[Shape[string], string, string]{}, sizePlusAlignOf[string](), ptrSize},
		{"result<bool, string>", &SizedResult[Shape[string], bool, string]{}, sizePlusAlignOf[string](), ptrSize},
		{"result<string, _>", &SizedResult[Shape[string], string, struct{}]{}, sizePlusAlignOf[string](), ptrSize},
		{"result<_, string>", &SizedResult[Shape[string], struct{}, string]{}, sizePlusAlignOf[string](), ptrSize},
		{"result<u64, u64>", &SizedResult[Shape[uint64], uint64, uint64]{}, 16, alignOf[uint64]()},
		{"result<u32, u64>", &SizedResult[Shape[uint64], uint32, uint64]{}, 16, alignOf[uint64]()},
		{"result<u64, u32>", &SizedResult[Shape[uint64], uint64, uint32]{}, 16, alignOf[uint64]()},
		{"result<u8, u64>", &SizedResult[Shape[uint64], uint8, uint64]{}, 16, alignOf[uint64]()},
		{"result<u64, u8>", &SizedResult[Shape[uint64], uint64, uint8]{}, 16, alignOf[uint64]()},
		{"result<u8, u32>", &SizedResult[Shape[uint32], uint8, uint32]{}, 8, alignOf[uint32]()},
		{"result<u32, u8>", &SizedResult[Shape[uint32], uint32, uint8]{}, 8, alignOf[uint32]()},
		{"result<[9]byte, u64>", &SizedResult[Shape[[9]byte], [9]byte, uint64]{}, 24, alignOf[uint64]()},

		{"result<string, _>", &OKSizedResult[string, struct{}]{}, sizePlusAlignOf[string](), ptrSize},
		{"result<string, bool>", &OKSizedResult[string, bool]{}, sizePlusAlignOf[string](), ptrSize},
		{"result<[9]byte, u64>", &OKSizedResult[[9]byte, uint64]{}, 24, alignOf[uint64]()},

		{"result<_, string>", &ErrSizedResult[struct{}, string]{}, sizePlusAlignOf[string](), ptrSize},
		{"result<bool, string>", &ErrSizedResult[bool, string]{}, sizePlusAlignOf[string](), ptrSize},
		{"result<u64, [9]byte>", &ErrSizedResult[uint64, [9]byte]{}, 24, alignOf[uint64]()},
	}

	for _, tt := range tests {
		typ := typeName(tt.r)
		t.Run(tt.name, func(t *testing.T) {
			if got, want := tt.r.Size(), tt.size; got != want {
				t.Errorf("(%s).Size() == %v, expected %v", typ, got, want)
			}
			if got, want := tt.r.ValOffset(), tt.offset; got != want {
				t.Errorf("(%s).ValOffset() == %v, expected %v", typ, got, want)
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
