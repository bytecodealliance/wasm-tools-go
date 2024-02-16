package cm

// TODO: should there be separate types for owned vs borrowed handles?
type Resource uint32

// ResourceNone is a sentinel value indicating a null or uninitialized resource.
// This is a reserved value specified in the [Canonical ABI runtime state].
//
// [Canonical ABI runtime state]: https://github.com/WebAssembly/component-model/blob/main/design/mvp/CanonicalABI.md#runtime-state
const ResourceNone = 0

type Handle uint32

type Own[T any] Handle

type Borrow[T any] Handle
