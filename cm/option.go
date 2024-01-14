package cm

// Option represents a Component Model option<T> type.
// The first byte is a bool representing none or some,
// followed by storage for the associated type T.
type Option[T any] struct {
	isSome bool
	v      T
}

// None returns an Option[T] representing the none value,
// equivalent to the zero value of Option[T].
func None[T any]() Option[T] {
	return Option[T]{}
}

// Some returns an Option[T] representing the some value.
func Some[T any](v T) Option[T] {
	return Option[T]{
		isSome: true,
		v:      v,
	}
}

// None returns true if o represents the none value.
// None returns false if o represents the some value.
func (o Option[T]) None() bool {
	return !o.isSome
}

// Some returns T, true if o represents the some value.
// Some returns T, false if o represents the none value.
func (o Option[T]) Some() (T, bool) {
	return o.v, o.isSome
}
