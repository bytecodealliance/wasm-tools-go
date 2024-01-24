package cm

import (
	"testing"
	"unsafe"
)

var (
	_ result[string, bool] = &Result[string, string, bool]{}
	_ result[string, bool] = &OKSizedResult[string, bool]{}
	_ result[bool, string] = &ErrSizedResult[bool, string]{}
)

type result[OK, Err any] interface {
	IsErr() bool
	SetOK(OK)
	SetErr(Err)
	OK() (ok OK, isOK bool)
	Err() (err Err, isErr bool)
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
		{"result", UntypedResult(false), 1, 0},
		{"ok", UntypedResult(ResultOK), 1, 0},
		{"err", UntypedResult(ResultErr), 1, 0},

		{"result<_, _>", UnsizedResult[struct{}, struct{}](false), 1, 0},
		{"result<[0]u8, _>", UnsizedResult[[0]byte, struct{}](false), 1, 0},

		{"result<string, string>", Result[string, string, string]{}, sizePlusAlignOf[string](), ptrSize},
		{"result<bool, string>", Result[string, bool, string]{}, sizePlusAlignOf[string](), ptrSize},
		{"result<string, _>", Result[string, string, struct{}]{}, sizePlusAlignOf[string](), ptrSize},
		{"result<_, string>", Result[string, struct{}, string]{}, sizePlusAlignOf[string](), ptrSize},
		{"result<u64, u64>", Result[uint64, uint64, uint64]{}, 16, alignOf[uint64]()},
		{"result<u32, u64>", Result[uint64, uint32, uint64]{}, 16, alignOf[uint64]()},
		{"result<u64, u32>", Result[uint64, uint64, uint32]{}, 16, alignOf[uint64]()},
		{"result<u8, u64>", Result[uint64, uint8, uint64]{}, 16, alignOf[uint64]()},
		{"result<u64, u8>", Result[uint64, uint64, uint8]{}, 16, alignOf[uint64]()},
		{"result<u8, u32>", Result[uint32, uint8, uint32]{}, 8, alignOf[uint32]()},
		{"result<u32, u8>", Result[uint32, uint32, uint8]{}, 8, alignOf[uint32]()},
		{"result<[9]u8, u64>", Result[Shape[[9]byte], [9]byte, uint64]{}, 24, alignOf[uint64]()},

		{"result<string, _>", OKResult[string]{}, sizePlusAlignOf[string](), ptrSize},
		{"result<string, _>", OKSizedResult[string, struct{}]{}, sizePlusAlignOf[string](), ptrSize},
		{"result<string, bool>", OKSizedResult[string, bool]{}, sizePlusAlignOf[string](), ptrSize},
		{"result<[9]u8, u64>", OKSizedResult[[9]byte, uint64]{}, 24, alignOf[uint64]()},

		{"result<_, string>", ErrResult[string]{}, sizePlusAlignOf[string](), ptrSize},
		{"result<_, string>", ErrSizedResult[struct{}, string]{}, sizePlusAlignOf[string](), ptrSize},
		{"result<bool, string>", ErrSizedResult[bool, string]{}, sizePlusAlignOf[string](), ptrSize},
		{"result<u64, [9]u8>", ErrSizedResult[uint64, [9]byte]{}, 24, alignOf[uint64]()},
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
