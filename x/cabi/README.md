### Usage

The `cabi` package contains a single exported WebAssembly function `cabi_realloc` ([Canonical ABI] realloc). To use, import this package with `_`:

```
import _ "github.com/bytecodealliance/wasm-tools-go/cabi"
```

`cabi_realloc` is a WebAssembly [core function](https://www.w3.org/TR/wasm-core-2/syntax/modules.html#functions) that is validated to have the following core function type:

```
(func (param $originalPtr i32)
      (param $originalSize i32)
      (param $alignment i32)
      (param $newSize i32)
      (result i32))
```

The [Canonical ABI] will use realloc both to allocate (passing 0 for the first two parameters) and reallocate. If the Canonical ABI needs realloc, validation requires this option to be present (there is no default).

[Canonical ABI]: https://github.com/WebAssembly/component-model/blob/main/design/mvp/CanonicalABI.md
