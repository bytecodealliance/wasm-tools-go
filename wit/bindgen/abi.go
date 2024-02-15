package bindgen

import (
	"slices"

	"github.com/ydnar/wasm-tools-go/wit"
)

// variantShape returns the largest associated type in v.
// If v has multiple types with the same size, it returns the
// type that contains a pointer.
func variantShape(v *wit.Variant) wit.Type {
	types := v.Types()
	if len(types) == 0 {
		return nil
	}
	slices.SortFunc(types, func(a, b wit.Type) int {
		switch {
		case a.Size() > b.Size():
			return -1
		case a.Size() < b.Size():
			return 1
		case a.HasPointer() && !b.HasPointer():
			return -1
		case !a.HasPointer() && b.HasPointer():
			return 1
		default:
			return 0
		}
	})
	return types[0]
}

// variantAlign returns the type with the highest align value in v.
func variantAlign(v *wit.Variant) wit.Type {
	types := v.Types()
	if len(types) == 0 {
		return nil
	}
	slices.SortFunc(types, func(a, b wit.Type) int {
		switch {
		case a.Align() > b.Align():
			return -1
		case a.Align() < b.Align():
			return 1
		default:
			return 0
		}
	})
	return types[0]
}
