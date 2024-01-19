//go:build !cmboundscheck

package cm

// BoundsCheck determines if bounds-checking is compiled into this package.
// To enable bounds-checking, build this package with -tags cmboundscheck.
// The default is off.
const BoundsCheck = false
