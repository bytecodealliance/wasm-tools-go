# Design

This directory contains design notes and prototypes of Go bindings for Component Model interfaces.

## Exports

### Export Bindings for Resource Types

For each exported resource type, in addition to constructors, methods, and static functions, a component must import and export a number of functions to manage the resource lifecycle.

### Example

For a hypothetical WIT resource `water` in package/interface `example:foo/bar`, a component _imports_ the following administrative functions:

- **resource-new**: initialize a resource handle for a given `rep` (concrete representation). Called by the component before a resource value is returned by an exported function.
	- Import: `[export]example:foo/bar [resource-new]water`
	- Params: `i32` rep (in C bindings, this is a pointer to the representation)
	- Results: `i32` handle
- **resource-rep**: return the underlying representation of a handle. Can only be called by the component instance that created the original resource.
	- Import: `[export]example:foo/bar [resource-rep]water`
	- Params: `i32` handle
	- Results: `i32` rep (in C bindings, this is a pointer to the representation)
- **resource-drop**: drops an owned handle to a resource. If the resource has zero loans, the destructor (below) will be called.
	- Import: `example:foo/bar [resource-drop]water`
	- Params: `i32` handle
	- Results: none

In addition, the resource destructor is _exported_:

- **dtor**: destructor, implemented by user code
	- Export: `example:foo/bar#[dtor]water`
	- Params: `i32` rep (in C bindings, this is a pointer)
	- Results: none
	- Notes: the

Similarly, a function is exported for each resource method (e.g. `drink` and `spill`)

- Export: `example:foo/bar#[method]water.drink` (implemented by user code)
- Export: `example:foo/bar#[method]water.spill` (implemented by user code)

### Example in Go

#### Generated Bindings

```go
package bar

// imports omitted

type Water cm.Resource

//go:wasmimport [export]example:foo/bar [resource-new]water
func wasmimport_WaterResourceNew(rep uintptr) Water

func WaterResourceNew[Rep WaterInterface](rep Rep) Water {
	return wasmimport_WaterResourceNew(*(*uintptr)(unsafe.Pointer(&rep)))
}
```
