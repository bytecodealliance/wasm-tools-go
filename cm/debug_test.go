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
	DataAlign() uintptr
	DataOffset() uintptr
}

func (v Variant[Disc, Shape, Align]) Size() uintptr       { return unsafe.Sizeof(v) }
func (v Variant[Disc, Shape, Align]) DataAlign() uintptr  { return unsafe.Alignof(v.data) }
func (v Variant[Disc, Shape, Align]) DataOffset() uintptr { return tinyunsafe.OffsetOf(&v, &v.data) }

// ResultDebug is an interface used in tests to validate layout of result types.
type ResultDebug interface {
	VariantDebug
}

func (r Result) Size() uintptr       { return unsafe.Sizeof(r) }
func (r Result) DataAlign() uintptr  { return 0 }
func (r Result) DataOffset() uintptr { return 0 }

func (r OKResult[OK, Err]) Size() uintptr       { return unsafe.Sizeof(r) }
func (r OKResult[OK, Err]) DataAlign() uintptr  { return unsafe.Alignof(r.data) }
func (r OKResult[OK, Err]) DataOffset() uintptr { return tinyunsafe.OffsetOf(&r, &r.data) }

func (r ErrResult[OK, Err]) Size() uintptr       { return unsafe.Sizeof(r) }
func (r ErrResult[OK, Err]) DataAlign() uintptr  { return unsafe.Alignof(r.data) }
func (r ErrResult[OK, Err]) DataOffset() uintptr { return tinyunsafe.OffsetOf(&r, &r.data) }
