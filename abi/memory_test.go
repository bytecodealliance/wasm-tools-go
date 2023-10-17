package abi

import (
	"fmt"
	"slices"
	"testing"
	"unsafe"
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

func TestRealloc(t *testing.T) {
	const threshold = 16
	tests := []struct {
		name    string
		ptr     unsafe.Pointer
		size    uintptr
		align   uintptr
		newsize uintptr
		want    unsafe.Pointer
	}{
		{"nil", nil, 0, 1, 0, nil},
		{"nil with align", nil, 0, 2, 0, nil},
		{"align to 2", up(1), 0, 2, 0, up(2)},
		{"align to 8", up(1), 0, 8, 0, up(8)},
		{"align to 8", up(1), 0, 8, 0, up(8)},
		{"alloc 100 bytes", nil, 0, 1, 100, unsafe.Pointer(unsafe.SliceData(make([]byte, 100)))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Realloc(tt.ptr, tt.size, tt.align, tt.newsize)
			if (uintptr(tt.want) < threshold && got != tt.want) || (uintptr(tt.want) >= threshold && uintptr(got) < threshold) {
				t.Errorf("Realloc(%d, %d, %d, %d): expected %d, got %d",
					tt.ptr, tt.size, tt.align, tt.newsize, tt.want, got)
			}
			if uintptr(got) < threshold {
				return // it didn’t allocate
			}
			if tt.ptr == nil {
				wants := unsafe.Slice((*byte)(tt.want), tt.newsize)
				gots := unsafe.Slice((*byte)(got), tt.newsize)
				if slices.Compare(wants, gots) != 0 {
					t.Errorf("expected %v, got %v", wants, gots)
				}
			}
		})
	}
}

// Appease vet, see https://github.com/golang/go/issues/58625
func up(p uintptr) unsafe.Pointer {
	return *(*unsafe.Pointer)(unsafe.Pointer(&p))
}
