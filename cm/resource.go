package cm

// TODO: should there be separate types for owned vs borrowed handles?
type Resource uint32

type Handle uint32

type Own[T any] Handle

type Borrow[T any] Handle
