package cm

// Resource represents an opaque Component Model [resource handle].
// It is represented in the [Canonical ABI] as an 32-bit integer.
//
// [resource handle]: https://github.com/WebAssembly/component-model/blob/main/design/mvp/Explainer.md#handle-types
// [Canonical ABI]: https://github.com/WebAssembly/component-model/blob/main/design/mvp/CanonicalABI.md
type Resource uint32

// ResourceNone is a sentinel value indicating a null or uninitialized resource.
// This is a reserved value specified in the [Canonical ABI runtime state].
//
// [Canonical ABI runtime state]: https://github.com/WebAssembly/component-model/blob/main/design/mvp/CanonicalABI.md#runtime-state
const ResourceNone = 0
