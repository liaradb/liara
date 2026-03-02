package btree

import (
	"slices"
	"testing"
)

func TestChain(t *testing.T) {
	t.Parallel()

	c := newChain()

	values := []int{1, 2, 3, 4, 5}
	for _, v := range values {
		c.append(v)
	}

	result := make([]int, 0, len(values))
	for _, i := range c.items() {
		result = append(result, i.(int))
	}

	slices.Reverse(values)
	if !slices.Equal(result, values) {
		t.Errorf("incorrect result: %v, expected: %v", result, values)
	}
}

func TestChain__Crabbing(t *testing.T) {
	t.Parallel()

	c := newChain()

	values := []int{1, 2, 3, 4, 5}
	for _, v := range values {
		c.append(v)
	}

	for i := range len(values) {
		c.setReleased(byte(i))

		result := make([]int, 0, len(values))
		for _, i := range c.unreleasedItems() {
			result = append(result, i.(int))
		}

		if !slices.Equal(result, values[i:]) {
			t.Errorf("incorrect result: %v, expected: %v", result, values[i:])
		}
	}
}
