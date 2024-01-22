package cm

import (
	"testing"
	"unsafe"
)

var (
	_ variant2[struct{}, struct{}] = zeroPtr[UnsizedVariant2[struct{}, struct{}]]()
	_ variant2[struct{}, struct{}] = zeroPtr[UnsizedVariant2[struct{}, struct{}]]()
	_ variant2[string, bool]       = zeroPtr[Variant2[string, string, bool]]()
	_ variant2[bool, string]       = zeroPtr[Variant2[string, bool, string]]()
)

type variant2[T0, T1 any] interface {
	Tag() bool
	Case0() (T0, bool)
	Case1() (T1, bool)
	Set0(T0)
	Set1(T1)
}

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
		{"variant { string; string }", Variant2[string, string, string]{}, sizePlusAlignOf[string](), ptrSize},
		{"variant { bool; string }", Variant2[string, bool, string]{}, sizePlusAlignOf[string](), ptrSize},
		{"variant { string; _ }", Variant2[string, string, struct{}]{}, sizePlusAlignOf[string](), ptrSize},
		{"variant { _; _ }", Variant2[string, struct{}, string]{}, sizePlusAlignOf[string](), ptrSize},
		{"variant { u64; u64 }", Variant2[uint64, uint64, uint64]{}, 16, alignOf[uint64]()},
		{"variant { u32; u64 }", Variant2[uint64, uint32, uint64]{}, 16, alignOf[uint64]()},
		{"variant { u64; u32 }", Variant2[uint64, uint64, uint32]{}, 16, alignOf[uint64]()},
		{"variant { u8; u64 }", Variant2[uint64, uint8, uint64]{}, 16, alignOf[uint64]()},
		{"variant { u64; u8 }", Variant2[uint64, uint64, uint8]{}, 16, alignOf[uint64]()},
		{"variant { u8; u32 }", Variant2[uint32, uint8, uint32]{}, 8, alignOf[uint32]()},
		{"variant { u32; u8 }", Variant2[uint32, uint32, uint8]{}, 8, alignOf[uint32]()},
		{"variant { [9]u8, u64 }", Variant2[Shape[[9]byte], [9]byte, uint64]{}, 24, alignOf[uint64]()},
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

func TestAsBoundsCheck(t *testing.T) {
	if !BoundsCheck {
		return // TinyGo does not support t.SkipNow
	}
	defer func() {
		if recover() == nil {
			t.Errorf("As did not panic")
		}
	}()
	var v Variant[uint8, uint8, uint8]
	_ = As[string](&v)
}

func TestGetBoundsCheck(t *testing.T) {
	if !BoundsCheck {
		return // TinyGo does not support t.SkipNow
	}
	defer func() {
		if recover() == nil {
			t.Errorf("Get did not panic")
		}
	}()
	var v Variant[uint8, uint8, uint8]
	_, _ = Get[string](&v, 0)
}

func TestSetBoundsCheck(t *testing.T) {
	if !BoundsCheck {
		return // TinyGo does not support t.SkipNow
	}
	defer func() {
		if recover() == nil {
			t.Errorf("Set did not panic")
		}
	}()
	var v Variant[uint8, uint8, uint8]
	Set(&v, 0, "hello world")
}

func TestNewVariantBoundsCheck(t *testing.T) {
	if !BoundsCheck {
		return // TinyGo does not support t.SkipNow
	}
	defer func() {
		if recover() == nil {
			t.Errorf("NewVariant did not panic")
		}
	}()
	_ = NewVariant[uint8, uint8, uint8](0, "hello world")
}
