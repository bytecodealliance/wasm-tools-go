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
	if got, want := o1.Value(), (string)(""); got != want {
		t.Errorf("o1.Value: %v, expected %v", got, want)
	}

	var o2 Option[uint32]
	if got, want := o2.None(), true; got != want {
		t.Errorf("o2.None: %t, expected %t", got, want)
	}
	if got, want := o2.Some(), (*uint32)(nil); got != want {
		t.Errorf("o2.Some: %v, expected %v", got, want)
	}
	if got, want := o2.Value(), (uint32)(0); got != want {
		t.Errorf("o2.Value: %v, expected %v", got, want)
	}

	o3 := Some(true)
	if got, want := o3.None(), false; got != want {
		t.Errorf("o3.None: %t, expected %t", got, want)
	}
	if got, want := o3.Some(), &o3.some; got != want {
		t.Errorf("o3.Some: %v, expected %v", got, want)
	}
	if got, want := o3.Value(), true; got != want {
		t.Errorf("o3.Value: %v, expected %v", got, want)
	}

	o4 := Some("hello")
	if got, want := o4.None(), false; got != want {
		t.Errorf("o4.None: %t, expected %t", got, want)
	}
	if got, want := o4.Some(), &o4.some; got != want {
		t.Errorf("o4.Some: %v, expected %v", got, want)
	}
	if got, want := o4.Value(), "hello"; got != want {
		t.Errorf("o4.Value: %v, expected %v", got, want)
	}

	o5 := Some(List[string]{})
	if got, want := o5.None(), false; got != want {
		t.Errorf("o5.None: %t, expected %t", got, want)
	}
	if got, want := o5.Some(), &o5.some; got != want {
		t.Errorf("o5.Some: %v, expected %v", got, want)
	}
	if got, want := o5.Value(), (List[string]{}); got != want {
		t.Errorf("o4.Value: %v, expected %v", got, want)
	}

	f := func(s string) Option[string] {
		return Some(s)
	}
	if got, want := f("hello").Value(), "hello"; got != want {
		t.Errorf("Value: %v, expected %v", got, want)
	}
}
