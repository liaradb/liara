package slice

import "testing"

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
}
