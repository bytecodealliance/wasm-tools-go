// Package monotonicclock represents the interface "wasi:clocks/monotonic-clock".
//
// WASI Monotonic Clock is a clock API intended to let users measure elapsed
// time.
//
// It is intended to be portable at least between Unix-family platforms and
// Windows.
//
// A monotonic clock is a clock which has an unspecified initial value, and
// successive reads of the clock will produce non-decreasing values.
//
// It is intended for measuring elapsed time.
package monotonicclock

import "internal/wasm/wasi/io/poll"

// Instant represents the type "wasi:clocks/monotonic-clock.instant".
//
// An instant in time, in nanoseconds. An instant is relative to an
// unspecified initial value, and can only be compared to instances from
// the same monotonic-clock.
type Instant uint64

// Duration represents the type "wasi:clocks/monotonic-clock.duration".
//
// A duration of time, in nanoseconds.
type Duration uint64

// Now calls the imported function "wasi:clocks/monotonic-clock#now".
//
// Read the current value of the clock.
//
// The clock is monotonic, therefore calling this function repeatedly will
// produce a sequence of non-decreasing values.
func Now() Instant {
	return now()
}

//go:wasmimport wasi:clocks/monotonic-clock@0.2.0-rc-2023-11-10 now
func now() Instant

// Resolution calls the imported function "wasi:clocks/monotonic-clock#resolution".
//
// Query the resolution of the clock. Returns the duration of time
// corresponding to a clock tick.//
func Resolution() Instant {
	return resolution()
}

//go:wasmimport wasi:clocks/monotonic-clock@0.2.0-rc-2023-11-10 now
func resolution() Instant

// SubscribeInstant calls the imported function "wasi:clocks/monotonic-clock#subscribe-instant".
func SubscribeInstant(when Instant) poll.Pollable {
	return subscribe_instant(when)
}

//go:wasmimport wasi:clocks/monotonic-clock@0.2.0-rc-2023-11-10 subscribe-instant
func subscribe_instant(when Instant) poll.Pollable

// SubscribeDuration represents the imported function "wasi:clocks/monotonic-clock#subscribe-duration".
//
// Create a `pollable` which will resolve once the given duration has
// elapsed, starting at the time at which this function was called.
// occured.
func SubscribeDuration(when Duration) poll.Pollable {
	return subscribe_duration(when)
}

//go:wasmimport wasi:clocks/monotonic-clock@0.2.0-rc-2023-11-10 subscribe-duration
func subscribe_duration(when Duration) poll.Pollable
