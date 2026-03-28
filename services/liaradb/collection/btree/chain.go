package btree

import (
	"container/list"
	"iter"

	"github.com/liaradb/liaradb/collection/btree/keynode"
	"github.com/liaradb/liaradb/collection/btree/leafnode"
)

// TODO: Create latch crabbing support
type chain struct {
	l        *list.List
	released byte
}

func newChain() *chain {
	return &chain{l: list.New()}
}

func (c *chain) append(v any) {
	c.l.PushFront(v)
}

func (c *chain) setReleased(i byte) {
	c.released = i
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

func (c *chain) unreleasedItems() iter.Seq2[bool, any] {
	i := 0
	e := c.l.Back()
	iLeaf := c.l.Len() - 1
	return func(yield func(bool, any) bool) {
		for {
			if e == nil {
				return
			}
			if i >= int(c.released) {
				if !yield(i == iLeaf, e.Value) {
					return
				}
			}
			e = e.Prev()
			i++
		}
	}
}

func (c *chain) release() {
	for leaf, n := range c.unreleasedItems() {
		if leaf {
			ln := n.(*leafnode.LeafNode)
			ln.Release()
		} else {
			kn := n.(*keynode.KeyNode)
			kn.Release()
		}
	}
	c.setReleased(byte(c.l.Len()))
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
