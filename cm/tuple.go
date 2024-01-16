package cm

// Tuple represents a [Component Model tuple] with 2 fields.
//
// [Component Model tuple]: https://component-model.bytecodealliance.org/design/wit.html#tuples
type Tuple[T0, T1 any] struct {
	V0 T0
	V1 T1
}

// Tuple3 represents a [Component Model tuple] with 3 fields.
//
// [Component Model tuple]: https://component-model.bytecodealliance.org/design/wit.html#tuples
type Tuple3[T0, T1, T2 any] struct {
	V0 T0
	V1 T1
	V2 T2
}

// Tuple4 represents a [Component Model tuple] with 4 fields.
//
// [Component Model tuple]: https://component-model.bytecodealliance.org/design/wit.html#tuples
type Tuple4[T0, T1, T2, T3 any] struct {
	V0 T0
	V1 T1
	V2 T2
	V3 T3
}
