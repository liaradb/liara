package btreememory

import (
	"cmp"
	"slices"
)

type leafNode[K cmp.Ordered, V any] struct {
	k        K
	children []*leafEntry[K, V]
	parent   *keyNode[K, V]
	left     *leafNode[K, V]
	right    *leafNode[K, V]
}

var _ node[int, int] = (*leafNode[int, int])(nil)

func newLeafNode[K cmp.Ordered, V any](k K, v V) *leafNode[K, V] {
	return &leafNode[K, V]{
		k:        k,
		children: []*leafEntry[K, V]{newLeafEntry(k, v)},
	}
}

func (ln *leafNode[K, V]) key() K {
	return ln.k
}

func (ln *leafNode[K, V]) count() int {
	count := 0
	for _, l := range ln.children {
		count += l.count()
	}
	return count
}

func (ln *leafNode[K, V]) setParent(p *keyNode[K, V]) {
	ln.parent = p
}

func (ln *leafNode[K, V]) getValue(k K) (V, bool) {
	if ln == nil {
		return ln.zero()
	}

	return ln.getChild(k).getValue()
}

func (ln *leafNode[K, V]) getChild(k K) *leafEntry[K, V] {
	for _, l := range ln.children {
		if l.key == k {
			return l
		}
	}

	return nil
}

func (ln *leafNode[K, V]) insert(f int, k K, v V) (node[K, V], bool) {
	c := ln.getChild(k)
	if c != nil {
		// TODO: Create Overflow
		c.append(v)
		return nil, false
	}

	i := ln.getInsertionIndex(k)
	if i == 0 {
		ln.k = k
	}

	// TODO: Split before inserting
	ln.children = slices.Insert(ln.children, i, newLeafEntry(k, v))
	if len(ln.children) <= f {
		return nil, false
	}

	return ln.split(), true
}

func (ln *leafNode[K, V]) getInsertionIndex(k K) int {
	for i := len(ln.children) - 1; i >= 0; i-- {
		j := ln.children[i]
		if k >= j.key {
			return i + 1
		}
	}
	return 0
}

func (ln *leafNode[K, V]) split() node[K, V] {
	half := len(ln.children) / 2

	ln2 := &leafNode[K, V]{
		k:        ln.children[half].key,
		children: ln.children[half:],
		parent:   ln.parent,
		left:     ln,
		right:    ln.right,
	}

	// TODO: Should we copy slices?
	ln.children = slices.Clone(ln.children[:half])
	ln.right = ln2

	return ln2
}

func (ln *leafNode[K, V]) delete(f int, k K, v V) {

}

func (ln *leafNode[K, V]) deleteAll(f int, k K) {
	c, i := ln.getChildForDeletion(k)
	if c == nil {
		return
	}

	if ln.isMinimum(f) {
		// TODO: Rebalance
		if ln.left != nil && !ln.left.isMinimum(f) {
			// Borrow Left
			e := ln.left.popSmallest()
			// TODO: How do we handle overflow?
			ln.insert(f, e.key, e.value[0])
			// Pull smallest from Left
			// Update This Key
			// Key change propagates
		} else if ln.right != nil && !ln.right.isMinimum(f) {
			// Borrow Right
			e := ln.right.popLargest()
			// TODO: How do we handle overflow?
			ln.insert(f, e.key, e.value[0])
			// Pull largest from Right
			// Update Right Key
			// Key changes propagates
		} else if ln.left != nil {
			// Merge Left
			// Move children to Left
			// Delete node
			// Deletion propagates
		} else if ln.right != nil {
			// Merge Right
			// Move children to Right
			// Update Right Key
			// Key change propagates
			// Delete node
			// Deletion propagates
		} else {
			// Delete
		}
	} else {
		// Delete
	}
	ln.children = slices.Delete(ln.children, i, i+1)
}

func (ln *leafNode[K, V]) popLargest() *leafEntry[K, V] {
	largest := ln.children[0]
	ln.children = ln.children[1:]
	return largest
}

func (ln *leafNode[K, V]) popSmallest() *leafEntry[K, V] {
	i := len(ln.children) - 1
	smallest := ln.children[i]
	ln.children = ln.children[:i]
	return smallest
}

func (ln *leafNode[K, V]) getChildForDeletion(k K) (*leafEntry[K, V], int) {
	for i, l := range ln.children {
		if l.key == k {
			return l, i
		}
	}

	return nil, 0
}

func (ln *leafNode[K, V]) isMinimum(f int) bool {
	return len(ln.children) <= ln.minimum(f)
}

// TODO: Can we store this?
func (ln *leafNode[K, V]) minimum(f int) int {
	return ceiling(f, 2) - 1
}

func ceiling(a, b int) int {
	return (a + b - 1) / b
}

func (ln *leafNode[K, V]) height() int {
	if ln == nil {
		return 0
	}

	return 1
}

func (*leafNode[K, V]) zero() (V, bool) {
	var v V
	return v, false
}
