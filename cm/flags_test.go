package cm

import (
	"fmt"
	"testing"
)

func TestFlags(t *testing.T) {
	type MyFlag Flag
	const (
		F0 MyFlag = iota
		F1

		F32 MyFlag = 32
		F63 MyFlag = 63
		F95 MyFlag = 95
	)

	var flags1 struct {
		Flags[[2]uint32, MyFlag]
	}
	flags1.Set(F0)
	flags1.Set(63)
	if !flags1.IsSet(F0) {
		t.Errorf("expected bit 0 to be set")
	}
	if !flags1.IsSet(F63) {
		t.Errorf("expected bit %d to be set", F63)
	}
	if flags1.IsSet(F32) {
		t.Errorf("expected bit %d to not be set", F32)
	}
	fmt.Printf("flags1: %b\n", flags1.data)

	var flags2 struct {
		Flags[[3]uint32, MyFlag]
	}
	flags2.Set(F1)
	flags2.Set(F95)
	if !flags2.IsSet(F1) {
		t.Errorf("expected bit %d to be set", F1)
	}
	if !flags2.IsSet(F95) {
		t.Errorf("expected bit %d to be set", F95)
	}
	if flags2.IsSet(0) {
		t.Errorf("expected bit %d to not be set", F0)
	}
	fmt.Printf("flags2: %b\n", flags2.data)
}
