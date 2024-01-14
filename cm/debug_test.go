package cm

import (
	"reflect"
	"strings"
	"unsafe"

	"github.com/ydnar/wasm-tools-go/internal/tinyunsafe"
)

func typeName(v any) string {
	var name string
	if t := reflect.TypeOf(v); t.Kind() == reflect.Ptr {
		name = "*" + t.Elem().String()
	} else {
		name = t.String()
	}
	return strings.ReplaceAll(name, " ", "")
}

func sizePlusAlignOf[T any]() uintptr {
	var v T
	return unsafe.Sizeof(v) + unsafe.Alignof(v)
}

func alignOf[T any]() uintptr {
	var v T
	return unsafe.Alignof(v)
}

func zeroPtr[T any]() *T {
	var zero T
	return &zero
}

// VariantDebug is an interface used in tests to validate layout of variant types.
type VariantDebug interface {
	Size() uintptr
	ValAlign() uintptr
	ValOffset() uintptr
}

func (v SizedVariant2[Shape, T0, T1]) Size() uintptr {
	return unsafe.Sizeof(v)
}

func (v SizedVariant2[Shape, T0, T1]) ValAlign() uintptr {
	return unsafe.Alignof(v.val)
}

func (v SizedVariant2[Shape, T0, T1]) ValOffset() uintptr {
	return tinyunsafe.OffsetOf(&v, &v.val)
}

func (v UnsizedVariant2[T0, T1]) Size() uintptr {
	return unsafe.Sizeof(v)
}

func (v UnsizedVariant2[T0, T1]) ValAlign() uintptr {
	return 0
}

func (v UnsizedVariant2[T0, T1]) ValOffset() uintptr {
	return 0
}

// ResultDebug is an interface used in tests to validate layout of result types.
type ResultDebug interface {
	VariantDebug
}

func (r SizedResult[S, OK, Err]) Size() uintptr {
	return unsafe.Sizeof(r)
}

func (r SizedResult[S, OK, Err]) ValAlign() uintptr {
	return r.v.ValAlign()
}

func (r SizedResult[S, OK, Err]) ValOffset() uintptr {
	return r.v.ValOffset()
}

func (r UnsizedResult[OK, Err]) Size() uintptr {
	return unsafe.Sizeof(r)
}

func (r UnsizedResult[OK, Err]) ValAlign() uintptr {
	return r.v.ValAlign()
}

func (r UnsizedResult[OK, Err]) ValOffset() uintptr {
	return r.v.ValOffset()
}
