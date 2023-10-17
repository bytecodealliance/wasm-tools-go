package abi

import (
	"fmt"
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
	var sentinel uintptr
	sentinel -= 1 // wraparound
	tests := []struct {
		name    string
		ptr     uintptr
		size    uintptr
		align   uintptr
		newsize uintptr
		want    uintptr
	}{
		{"nil", 0, 0, 1, 0, 0},
		{"nil with align", 0, 0, 2, 0, 0},
		{"align to 2", 1, 0, 2, 0, 2},
		{"align to 8", 1, 0, 8, 0, 8},
		{"align to 8", 1, 0, 8, 0, 8},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Realloc(unsafePointer(tt.ptr), tt.size, tt.align, tt.newsize)
			if got != unsafePointer(tt.want) {
				t.Errorf("Realloc(%d, %d, %d, %d): expected %d, got %d",
					tt.ptr, tt.size, tt.align, tt.newsize, tt.want, got)
			}
		})
	}
}

// Appease vet, see https://github.com/golang/go/issues/58625
func unsafePointer(p uintptr) unsafe.Pointer {
	return *(*unsafe.Pointer)(unsafe.Pointer(&p))
}
