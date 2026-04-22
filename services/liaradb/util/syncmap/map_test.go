package syncmap

import (
	"maps"
	"testing"
)

func TestMap_Default(t *testing.T) {
	t.Parallel()

	m := SyncMap[string, int]{}

	c := 0
	for range m.Iterate() {
		c++
	}

	if c != 0 {
		t.Errorf("incorrect count: %v, expected: %v", c, 0)
	}
}

func TestMap(t *testing.T) {
	t.Parallel()

	m := SyncMap[string, int]{}

	data := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}

	for k, v := range data {
		m.Set(k, v)
	}

	for k, v := range data {
		r, ok := m.Get(k)
		if !ok {
			t.Errorf("key should exist: %v", k)
		}
		if r != v {
			t.Errorf("incorrect result: %v, expected: %v", r, v)
		}
	}

	for range m.Iterate() {
		// early return
		break
	}

	result := maps.Collect(m.Iterate())

	if !maps.Equal(result, data) {
		t.Errorf("incorrect result: %v, expected: %v", result, data)
	}
}
