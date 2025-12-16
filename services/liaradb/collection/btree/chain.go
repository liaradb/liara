package btree

import (
	"container/list"
	"iter"
)

// TODO: Test this
type chain struct {
	l *list.List
}

func newChain(l *list.List) *chain {
	return &chain{l: l}
}

func (c *chain) items() iter.Seq2[int, any] {
	i := 0
	return func(yield func(int, any) bool) {
		e := c.l.Front()
		if e == nil || !yield(i, e.Value) {
			return
		}
		i++
	}
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
