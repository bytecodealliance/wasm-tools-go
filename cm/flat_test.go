package cm

import "testing"

func TestLower(t *testing.T) {
	print("int32   ")
	_ = Lower[int32](0)
	print("uint32  ")
	_ = Lower[uint32](0)
	print("int64   ")
	_ = Lower[int64](0)
	print("uint64  ")
	_ = Lower[uint64](0)
	print("float32 ")
	_ = Lower[float32](0)
	print("float64 ")
	_ = Lower[float64](0)
}
