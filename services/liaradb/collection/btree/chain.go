package btree

import (
	"container/list"
	"iter"

	"github.com/liaradb/liaradb/collection/btree/keynode"
	"github.com/liaradb/liaradb/collection/btree/leafnode"
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
			ln := n.(*leafnode.LeafNode)
			ln.Release()
		} else {
			kn := n.(*keynode.KeyNode)
			kn.Release()
		}
	}
}

func (c *chain) latch() {
	for i, n := range c.items() {
		if i == 0 {
			ln := n.(*leafnode.LeafNode)
			ln.Latch()
		} else {
			kn := n.(*keynode.KeyNode)
			kn.Latch()
		}
	}
}

func (c *chain) unlatch() {
	for i, n := range c.items() {
		if i == 0 {
			ln := n.(*leafnode.LeafNode)
			ln.Unlatch()
		} else {
			kn := n.(*keynode.KeyNode)
			kn.Unlatch()
		}
	}
}
