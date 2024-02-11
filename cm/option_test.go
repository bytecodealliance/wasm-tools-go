package cm

import "testing"

func TestOption(t *testing.T) {
	o1 := None[string]()
	if got, want := o1.None(), true; got != want {
		t.Errorf("o1.None: %t, expected %t", got, want)
	}
	if got, want := o1.Some(), (*string)(nil); got != want {
		t.Errorf("o1.Some: %v, expected %v", got, want)
	}

	var o2 Option[uint32]
	if got, want := o2.None(), true; got != want {
		t.Errorf("o2.None: %t, expected %t", got, want)
	}
	if got, want := o2.Some(), (*uint32)(nil); got != want {
		t.Errorf("o2.Some: %v, expected %v", got, want)
	}

	o3 := Some(true)
	if got, want := o3.None(), false; got != want {
		t.Errorf("o3.None: %t, expected %t", got, want)
	}
	if got, want := o3.Some(), &o3.some; got != want {
		t.Errorf("o3.Some: %v, expected %v", got, want)
	}
}
