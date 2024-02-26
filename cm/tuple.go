package cm

// Tuple represents a [Component Model tuple] with 2 fields.
//
// [Component Model tuple]: https://component-model.bytecodealliance.org/design/wit.html#tuples
type Tuple[T0, T1 any] struct {
	F0 T0
	F1 T1
}

// Tuple3 represents a [Component Model tuple] with 3 fields.
//
// [Component Model tuple]: https://component-model.bytecodealliance.org/design/wit.html#tuples
type Tuple3[T0, T1, T2 any] struct {
	F0 T0
	F1 T1
	F2 T2
}

// Tuple4 represents a [Component Model tuple] with 4 fields.
//
// [Component Model tuple]: https://component-model.bytecodealliance.org/design/wit.html#tuples
type Tuple4[T0, T1, T2, T3 any] struct {
	F0 T0
	F1 T1
	F2 T2
	F3 T3
}

// Tuple5 represents a [Component Model tuple] with 5 fields.
//
// [Component Model tuple]: https://component-model.bytecodealliance.org/design/wit.html#tuples
type Tuple5[T0, T1, T2, T3, T4 any] struct {
	F0 T0
	F1 T1
	F2 T2
	F3 T3
	F4 T4
}

// Tuple6 represents a [Component Model tuple] with 6 fields.
//
// [Component Model tuple]: https://component-model.bytecodealliance.org/design/wit.html#tuples
type Tuple6[T0, T1, T2, T3, T4, T5 any] struct {
	F0 T0
	F1 T1
	F2 T2
	F3 T3
	F4 T4
	F5 T5
}

// Tuple7 represents a [Component Model tuple] with 7 fields.
//
// [Component Model tuple]: https://component-model.bytecodealliance.org/design/wit.html#tuples
type Tuple7[T0, T1, T2, T3, T4, T5, T6 any] struct {
	F0 T0
	F1 T1
	F2 T2
	F3 T3
	F4 T4
	F5 T5
	F6 T6
}

// Tuple8 represents a [Component Model tuple] with 8 fields.
//
// [Component Model tuple]: https://component-model.bytecodealliance.org/design/wit.html#tuples
type Tuple8[T0, T1, T2, T3, T4, T5, T6, T7 any] struct {
	F0 T0
	F1 T1
	F2 T2
	F3 T3
	F4 T4
	F5 T5
	F6 T6
	F7 T7
}

// MaxTuple specifies the maximum number of fields in a Tuple(n) type.
// currently [Tuple8].
const MaxTuple = 8
