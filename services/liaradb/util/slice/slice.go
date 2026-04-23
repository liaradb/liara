package slice

import "iter"

type Slice[V comparable] struct {
	values []V
}

func New[V comparable](size int) *Slice[V] {
	return &Slice[V]{
		values: make([]V, 0, size),
	}
}

func (c *Slice[V]) Length() int {
	return len(c.values)
}

func (c *Slice[V]) Get(i int) (V, bool) {
	if i < 0 || i >= len(c.values) {
		return c.zero()
	}

	return c.values[i], true
}

func (c *Slice[V]) Pop() (V, bool) {
	length := len(c.values)
	if length <= 0 {
		return c.zero()
	}

	last := length - 1
	v := c.values[last]
	c.values = c.values[:last]

	return v, true
}

func (c *Slice[V]) zero() (V, bool) {
	var v V
	return v, false
}

func (c *Slice[V]) Push(v V) {
	c.values = append(c.values, v)
}

func (c *Slice[V]) SwapWithLast(i int) {
	j := len(c.values) - 1
	if i < 0 || i >= j {
		return
	}

	c.values[i], c.values[j] = c.values[j], c.values[i]
}

func (c *Slice[V]) SwapAndPop(i int) (V, bool) {
	c.SwapWithLast(i)
	return c.Pop()
}

func (c *Slice[V]) Find(v V) (int, bool) {
	for i, j := range c.values {
		if j == v {
			return i, true
		}
	}

	return -1, false
}

func (c *Slice[V]) RemoveAnyOrder(v V) bool {
	i, ok := c.Find(v)
	if !ok {
		return false
	}

	_, ok = c.SwapAndPop(i)
	return ok
}

func (c *Slice[V]) Iterate() iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range c.values {
			if !yield(v) {
				return
			}
		}
	}
}
