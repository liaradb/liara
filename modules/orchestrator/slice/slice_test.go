package slice

import (
	"slices"
	"testing"
)

func TestSlice(t *testing.T) {
	t.Parallel()

	t.Run("should push and pop", func(t *testing.T) {
		t.Parallel()

		c := New[int](0)

		c.Push(1)
		c.Push(2)

		if v, ok := c.Pop(); !ok || v != 2 {
			t.Error("value is incorrect")
		}

		if v, ok := c.Pop(); !ok || v != 1 {
			t.Error("value is incorrect")
		}
	})

	t.Run("should not pop beyond empty", func(t *testing.T) {
		t.Parallel()

		c := New[int](0)

		if v, ok := c.Pop(); ok || v != 0 {
			t.Error("value is incorrect")
		}
	})

	t.Run("should get length", func(t *testing.T) {
		t.Parallel()

		c := New[int](0)

		c.Push(1)
		if c.Length() != 1 {
			t.Error("incorrect length")
		}

		c.Push(2)
		if c.Length() != 2 {
			t.Error("incorrect length")
		}
	})

	t.Run("should get", func(t *testing.T) {
		t.Parallel()

		c := New[int](0)

		c.Push(1)
		c.Push(2)

		if v, ok := c.Get(0); !ok || v != 1 {
			t.Error("value is incorrect")
		}

		if v, ok := c.Get(1); !ok || v != 2 {
			t.Error("value is incorrect")
		}
	})

	t.Run("should not get out of range", func(t *testing.T) {
		c := New[int](0)

		c.Push(1)
		c.Push(2)

		if v, ok := c.Get(-1); ok || v != 0 {
			t.Error("value is incorrect")
		}

		if v, ok := c.Get(2); ok || v != 0 {
			t.Error("value is incorrect")
		}
	})

	t.Run("should swap", func(t *testing.T) {
		t.Parallel()

		c := New[int](0)

		c.Push(1)
		c.Push(2)
		c.Push(3)

		c.SwapWithLast(0)

		if v, ok := c.Get(0); !ok || v != 3 {
			t.Error("value is incorrect")
		}

		if v, ok := c.Get(2); !ok || v != 1 {
			t.Error("value is incorrect")
		}
	})

	t.Run("should not swap outside range", func(t *testing.T) {
		t.Parallel()

		c := New[int](0)

		c.Push(1)
		c.Push(2)
		c.Push(3)

		c.SwapWithLast(-1)
		c.SwapWithLast(3)

		if v, ok := c.Get(0); !ok || v != 1 {
			t.Error("value is incorrect")
		}

		if v, ok := c.Get(1); !ok || v != 2 {
			t.Error("value is incorrect")
		}

		if v, ok := c.Get(2); !ok || v != 3 {
			t.Error("value is incorrect")
		}
	})

	t.Run("should swap and pop", func(t *testing.T) {
		t.Parallel()

		c := New[int](0)

		c.Push(1)
		c.Push(2)
		c.Push(3)

		c.SwapAndPop(0)

		if c.Length() != 2 {
			t.Error("length is incorrect")
		}

		if v, ok := c.Get(0); !ok || v != 3 {
			t.Error("value is incorrect")
		}

		if v, ok := c.Get(1); !ok || v != 2 {
			t.Error("value is incorrect")
		}
	})

	t.Run("should remove in any order", func(t *testing.T) {
		t.Parallel()

		c := New[int](0)

		c.Push(1)
		c.Push(2)
		c.Push(3)

		c.RemoveAnyOrder(1)

		if c.Length() != 2 {
			t.Error("length is incorrect")
		}

		if v, ok := c.Get(0); !ok || v != 3 {
			t.Error("value is incorrect")
		}

		if v, ok := c.Get(1); !ok || v != 2 {
			t.Error("value is incorrect")
		}
	})

	t.Run("should find", func(t *testing.T) {
		t.Parallel()

		c := New[int](0)

		c.Push(1)
		c.Push(2)
		c.Push(3)

		if i, ok := c.Find(1); !ok || i != 0 {
			t.Error("value is incorrect")
		}

		if i, ok := c.Find(2); !ok || i != 1 {
			t.Error("value is incorrect")
		}

		if i, ok := c.Find(3); !ok || i != 2 {
			t.Error("value is incorrect")
		}

		if i, ok := c.Find(4); ok || i != -1 {
			t.Error("value is incorrect")
		}
	})

	t.Run("should iterate", func(t *testing.T) {
		t.Parallel()

		c := New[int](0)

		c.Push(1)
		c.Push(2)
		c.Push(3)

		for range c.Iterate() {
			break
		}

		want := []int{1, 2, 3}
		result := slices.Collect(c.Iterate())
		if !slices.Equal(result, want) {
			t.Errorf("incorrect result: %v, expected: %v", result, want)
		}
	})

	t.Run("should remove", func(t *testing.T) {
		t.Parallel()

		c := New[int](0)

		c.Push(1)
		c.Push(2)
		c.Push(3)

		if !c.RemoveAnyOrder(2) {
			t.Error("should remove")
		}

		if c.RemoveAnyOrder(4) {
			t.Error("should not remove")
		}

		if i, ok := c.Find(1); !ok || i != 0 {
			t.Error("value is incorrect")
		}

		if i, ok := c.Find(2); ok || i != -1 {
			t.Error("value is incorrect")
		}

		if i, ok := c.Find(3); !ok || i != 1 {
			t.Error("value is incorrect")
		}

		if i, ok := c.Find(4); ok || i != -1 {
			t.Error("value is incorrect")
		}
	})
}
