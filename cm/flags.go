package cm

import "unsafe"

type (
	Flags8  uint8
	Flags16 uint16
	Flags32 uint32
	Flags64 uint64
	Flags96 [3]uint32
)

// TODO: necessary?
func (f Flags8) Is(flag Flags8) bool {
	return f&flag != 0
}

type Flags[Shape any] struct {
	data Shape
}

func (f *Flags[Shape]) IsSet(bit uint) bool {
	ptr := (*byte)(unsafe.Add(unsafe.Pointer(f), bit>>3))
	return *ptr&(1<<(bit&0b111)) != 0
}

func (f *Flags[Shape]) Set(bit uint) {
	ptr := (*byte)(unsafe.Add(unsafe.Pointer(f), bit>>3))
	*ptr |= 1 << (bit & 0b111)
}

func (f *Flags[Shape]) Clear(bit uint) {
	ptr := (*byte)(unsafe.Add(unsafe.Pointer(f), bit>>3))
	*ptr &^= 1 << (bit & 0b111)
}
