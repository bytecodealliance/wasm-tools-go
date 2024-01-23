package cm

import (
	"fmt"
	"testing"
)

func TestFlags(t *testing.T) {
	var flags1 Flags[[2]uint32]
	flags1.Set(0)
	flags1.Set(63)
	if !flags1.IsSet(0) {
		t.Errorf("expected bit 0 to be set")
	}
	if !flags1.IsSet(63) {
		t.Errorf("expected bit 63 to be set")
	}
	if flags1.IsSet(32) {
		t.Errorf("expected bit 32 to not be set")
	}
	fmt.Printf("flags1: %b\n", flags1.data)

	var flags2 Flags[[3]uint32]
	flags2.Set(1)
	flags2.Set(95)
	if !flags2.IsSet(1) {
		t.Errorf("expected bit 0 to be set")
	}
	if !flags2.IsSet(95) {
		t.Errorf("expected bit 95 to be set")
	}
	if flags2.IsSet(0) {
		t.Errorf("expected bit 0 to not be set")
	}
	fmt.Printf("flags2: %b\n", flags2.data)
}
