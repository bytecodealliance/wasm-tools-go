package cm

import (
	"math"
	"testing"
)

func TestTuple(t *testing.T) {
	var H HostLayout
	_ = Tuple[string, bool]{H, "hello", false}
	_ = Tuple3[string, bool, uint8]{H, "hello", false, 1}
	_ = Tuple4[string, bool, uint8, uint16]{H, "hello", false, 1, 32000}
	_ = Tuple5[string, bool, uint8, uint16, uint32]{H, "hello", false, 1, 32000, 1_000_000}
	_ = Tuple6[string, bool, uint8, uint16, uint32, uint64]{H, "hello", false, 1, 32000, 1_000_000, 5_000_000_000}
	_ = Tuple7[string, bool, uint8, uint16, uint32, uint64, float32]{H, "hello", false, math.MaxUint8, math.MaxUint16, math.MaxUint32, math.MaxUint64, math.MaxFloat32}
	_ = Tuple8[string, bool, uint8, uint16, uint32, uint64, float32, float64]{H, "hello", false, math.MaxUint8, math.MaxUint16, math.MaxUint32, math.MaxUint64, math.MaxFloat32, math.MaxFloat64}
}
