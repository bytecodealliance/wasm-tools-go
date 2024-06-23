package cm

// Option represents a Component Model [option<T>] type.
// The first byte is a bool representing none or some,
// followed by storage for the associated type T.
//
// [option<T>]: https://component-model.bytecodealliance.org/design/wit.html#options
type Option[T any] struct {
	isSome bool
	some   T
}

// None returns an [Option] representing the none case,
// equivalent to the zero value.
func None[T any]() Option[T] {
	return Option[T]{}
}

// Some returns an [Option] representing the some case.
func Some[T any](v T) Option[T] {
	return Option[T]{
		isSome: true,
		some:   v,
	}
}

// None returns true if o represents the none case.
func (o *Option[T]) None() bool {
	return !o.isSome
}

// Some returns a non-nil *T if o represents the some case,
// or nil if o represents the none case.
func (o *Option[T]) Some() *T {
	if o.isSome {
		return &o.some
	}
	return nil
}
