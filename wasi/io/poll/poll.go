//go:build wasip2

// Package poll represents the imported interface "wasi:io/poll".
//
// A poll API intended to let users wait for I/O events on multiple handles
// at once.
package poll

import "github.com/ydnar/wasm-tools-go/cm"

// Pollable represents the imported type "wasi:io/poll.pollable".
//
// `pollable` represents a single I/O event which may be ready, or not.
type Pollable cm.Resource

// Ready calls the imported method "wasi:io/poll.pollable#ready".
//
// Return the readiness of a pollable. This function never blocks.
//
// Returns `true` when the pollable is ready, and `false` otherwise.
func (self Pollable) Ready() bool {
	return self.ready()
}

//go:wasmimport wasi:io/poll@0.2.0-rc-2023-11-10 [method]pollable.ready
func (self Pollable) ready() bool

// Ready represents the imported method "wasi:io/poll.pollable#block".
//
// `block` returns immediately if the pollable is ready, and otherwise
// blocks until ready.
//
// This function is equivalent to calling `poll.poll` on a list
// containing only this pollable.
func (self Pollable) Block() {
	self.block()
}

//go:wasmimport wasi:io/poll@0.2.0-rc-2023-11-10 [method]pollable.block
func (self Pollable) block()

// Poll represents the imported function "wasi:io/poll#poll".
//
// Poll for completion on a set of pollables.
//
// This function takes a list of pollables, which identify I/O sources of
// interest, and waits until one or more of the events is ready for I/O.
//
// The result `list<u32>` contains one or more indices of handles in the
// argument list that is ready for I/O.
//
// If the list contains more elements than can be indexed with a `u32`
// value, this function traps.
//
// A timeout can be implemented by adding a pollable from the
// wasi-clocks API to the list.
//
// This function does not return a `result`; polling in itself does not
// do any I/O so it doesn't fail. If any of the I/O sources identified by
// the pollables has an error, it is indicated by marking the source as
// being reaedy for I/O.
func Poll(in cm.List[Pollable]) cm.List[uint32] {
	var ret cm.List[uint32]
	poll(in, &ret)
	return ret
}

//go:wasmimport wasi:io/poll@0.2.0-rc-2023-11-10 pollable.poll
func poll(in cm.List[Pollable], ret *cm.List[uint32])
