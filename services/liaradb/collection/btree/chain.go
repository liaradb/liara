package btree

import (
	"container/list"
	"iter"
)

// TODO: Test this
type chain struct {
	l *list.List
}

func newChain() *chain {
	return &chain{l: list.New()}
}

func (c *chain) items() iter.Seq2[int, any] {
	i := 0
	e := c.l.Front()
	return func(yield func(int, any) bool) {
		for {
			if e == nil || !yield(i, e.Value) {
				return
			}
			e = e.Next()
			i++
		}
	}
}

func (c *chain) append(v any) {
	c.l.PushFront(v)
}

func (c *chain) release() {
	for i, n := range c.items() {
		if i == 0 {
			ln := n.(*LeafNode)
			ln.page.Release()
		} else {
			kn := n.(*KeyNode)
			kn.page.Release()
		}
	}
}
