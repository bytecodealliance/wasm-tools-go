package cm

import "unsafe"

// Flag represents an individual flag. The value of a Flag
// is index of the bit into a Flags value.
//
// The intended use is as a separate named type:
//
//	type MyFlag cm.Flag
//	type MyFlags struct { cm.Flags[uint8, MyFlag] }
type Flag uint

// FIXME(ydnar): this is silly
type flagsShape interface {
	uint8 | uint16 | uint32 |
		[2]uint32 | [3]uint32 | [4]uint32 | [5]uint32 | [6]uint32 | [7]uint32 |
		[8]uint32 | [9]uint32 | [10]uint32 | [11]uint32 | [12]uint32 | [13]uint32 | [14]uint32 | [15]uint32 |
		[16]uint32 | [17]uint32 | [18]uint32 | [19]uint32 | [20]uint32 | [21]uint32 | [22]uint32 | [23]uint32 |
		[24]uint32 | [25]uint32 | [26]uint32 | [27]uint32 | [28]uint32 | [29]uint32 | [30]uint32 | [31]uint32
}

// Flags represents a bitfield of multiple flags.
// Shape must be either uint8, uint16, uint32, or an
// array of uint32 large enough to contain the max
// value of the associated Flag type.
type Flags[Shape flagsShape, Flag ~uint] struct {
	data Shape
}

// IsSet returns true if flag is set.
func (f *Flags[Shape, Flag]) IsSet(flag Flag) bool {
	ptr := (*byte)(unsafe.Add(unsafe.Pointer(f), flag>>3))
	return *ptr&(1<<(flag&0b111)) != 0
}

// Set sets the bit indexed by flag to 1.
func (f *Flags[Shape, Flag]) Set(flag Flag) {
	ptr := (*byte)(unsafe.Add(unsafe.Pointer(f), flag>>3))
	*ptr |= 1 << (flag & 0b111)
}

// Clear clears the bit indexed by flag.
func (f *Flags[Shape, Flag]) Clear(flag Flag) {
	ptr := (*byte)(unsafe.Add(unsafe.Pointer(f), flag>>3))
	*ptr &^= 1 << (flag & 0b111)
}
