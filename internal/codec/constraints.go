package codec

import "math"

// Signed is the set of signed integer types supported by this package.
type Signed interface {
	int | int8 | int16 | int32 | int64
}

// Unsigned is the set of unsigned integer types supported by this package.
type Unsigned interface {
	uint | uint8 | uint16 | uint32 | uint64
}

// Integer is the set of integer types supported by this package.
type Integer interface {
	Signed | Unsigned
}

// Float is the set of floating-point types supported by this package.
type Float interface {
	float32 | float64
}

// TypeName returns the Go type name for supported integer and floating-point types.
func TypeName[T Integer | Float]() string {
	var v T
	v -= 1

	if v > 0 {
		// Unsigned int
		switch uint64(v) {
		case math.MaxUint64:
			return "uint64"
		case math.MaxUint32:
			return "uint32"
		case math.MaxUint16:
			return "uint16"
		}
		return "uint8"
	}

	// Signed int or float
	var max uint64 = math.MaxInt64
	v = T(max)
	if uint64(v) == max {
		return "int64"
	}

	if v > 0 {
		// Must be a float
		var max float64 = math.MaxFloat64
		v = T(max)
		if float64(v) == max {
			return "float64"
		}
		return "float32"
	}

	max = math.MaxInt32
	v = T(max)
	if uint64(v) == max {
		return "int32"
	}

	max = math.MaxInt16
	v = T(max)
	if uint64(v) == max {
		return "int16"
	}

	return "int8"
}
