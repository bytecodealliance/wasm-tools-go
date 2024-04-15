package ordered

import (
	"testing"
)

func TestMap(t *testing.T) {
	m := New[int, int]()
	m.Set(0, 0)
	m.Set(5, 5)
	m.Set(1, 1)
	m.Delete(5)
	m.Set(2, 2)
	m.Set(3, 3)
	m.Set(3, 3)
	m.Set(4, 4)
	m.Set(5, 5)

	// Test values
	for i := 0; i < 5; i++ {
		got, want := m.Get(i), i
		if got != want {
			t.Errorf("m.Get(%d): %d, expected %d", i, got, want)
		}
	}

	// Test iteration
	i := 0
	m.All()(func(k int, v int) bool {
		if k != i {
			t.Errorf("m.All() @ %d: k == %d, expected %d", i, k, i)
		}
		if v != i {
			t.Errorf("m.All() @ %d: v == %d, expected %d", i, v, i)
		}
		i++
		return true
	})

	// Test appending items during iteration
	i = 0
	m.All()(func(k int, v int) bool {
		if i == 3 {
			m.Set(6, 6)
		}
		if k != i {
			t.Errorf("m.All() @ %d: k == %d, expected %d", i, k, i)
		}
		if v != i {
			t.Errorf("m.All() @ %d: v == %d, expected %d", i, v, i)
		}
		i++
		return true
	})
	if i != 7 {
		t.Errorf("i == %d, expected 7", i)
	}

	// Test deleting items during iteration
	i = 0
	m.All()(func(k int, v int) bool {
		if k != i {
			t.Errorf("m.All() @ %d: k == %d, expected %d", i, k, i)
		}
		if v != i {
			t.Errorf("m.All() @ %d: v == %d, expected %d", i, v, i)
		}
		i++
		m.Delete(k)
		return true
	})
}
