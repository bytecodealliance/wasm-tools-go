package cm

import (
	"unsafe"
)

// TODO: this file represents an experiment for how to represent Component Model
// flags types larger than a uint64, which according to the Canonical ABI are
// represented by [N]uint32 where N >= the number of unique flags represented
// by the type.

// Flag represents an individual flag. The value of a Flag
// is index of the bit into a Flags value.
//
// The intended use is as a separate named type:
//
//	type MyFlag cm.Flag
//	type MyFlags struct { cm.Flags[uint8, MyFlag] }
type Flag uint

// flagsShape defines sufficient shapes to store up to 1024 flag values.
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

// Is returns true if flag is set.
func (f *Flags[Shape, Flag]) Is(flag Flag) bool {
	ptr := (*byte)(unsafe.Add(unsafe.Pointer(f), flag>>3))
	return *ptr&(1<<(flag&0b111)) != 0
}

// Set sets or clears the bit indexed by flag.
func (f *Flags[Shape, Flag]) Set(flag Flag) {
	ptr := (*byte)(unsafe.Add(unsafe.Pointer(f), flag>>3))
	*ptr |= 1 << (flag & 0b111)
}

// Clear clears the bit indexed by flag.
func (f *Flags[Shape, Flag]) Clear(flag Flag) {
	ptr := (*byte)(unsafe.Add(unsafe.Pointer(f), flag>>3))
	*ptr &^= 1 << (flag & 0b111)
}

// Flags8 represents a flags value with 1-8 unique flags.
type Flags8[Flag ~uint] uint8

// Is returns true if flag is set.
func (f *Flags8[Flag]) Is(flag Flag) bool {
	return *f&(1<<flag) != 0
}

// Set sets or clears the bit indexed by flag.
func (f *Flags8[Flag]) Set(flag Flag) {
	*f |= 1 << flag
}

// Clear clears the bit indexed by flag.
func (f *Flags8[Flag]) Clear(flag Flag) {
	*f &^= 1 << flag
}

// Flags16 represents a flags value with 9-16 unique flags.
type Flags16[Flag ~uint] uint16

// Is returns true if flag is set.
func (f *Flags16[Flag]) Is(flag Flag) bool {
	return *f&(1<<flag) != 0
}

// Set sets or clears the bit indexed by flag.
func (f *Flags16[Flag]) Set(flag Flag) {
	*f |= 1 << flag
}

// Clear clears the bit indexed by flag.
func (f *Flags16[Flag]) Clear(flag Flag) {
	*f &^= 1 << flag
}

// Flags32 represents a flags value with 17-32 unique flags.
type Flags32[Flag ~uint] uint32

// Is returns true if flag is set.
func (f *Flags32[Flag]) Is(flag Flag) bool {
	return *f&(1<<flag) != 0
}

// Set sets or clears the bit indexed by flag.
func (f *Flags32[Flag]) Set(flag Flag) {
	*f |= 1 << flag
}

// Clear clears the bit indexed by flag.
func (f *Flags32[Flag]) Clear(flag Flag) {
	*f &^= 1 << flag
}

// Flags64 represents a flags value with 33-64 unique flags.
type Flags64[Flag ~uint] uint64

// Is returns true if flag is set.
func (f *Flags64[Flag]) Is(flag Flag) bool {
	return *f&(1<<flag) != 0
}

// Set sets or clears the bit indexed by flag.
func (f *Flags64[Flag]) Set(flag Flag) {
	*f |= 1 << flag
}

// Clear clears the bit indexed by flag.
func (f *Flags64[Flag]) Clear(flag Flag) {
	*f &^= 1 << flag
}
