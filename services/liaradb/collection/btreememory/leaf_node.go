package btreememory

import (
	"cmp"
	"slices"
)

type leafNode[K cmp.Ordered] struct {
	storage  Storage[K]
	k        K
	children []*leafEntry[K]
	left     *leafNode[K]
	right    *leafNode[K]
}

var _ node[int] = (*leafNode[int])(nil)

func newLeafNode[K cmp.Ordered](s Storage[K], k K, rid RecordID) *leafNode[K] {
	return &leafNode[K]{
		storage:  s,
		k:        k,
		children: []*leafEntry[K]{newLeafEntry(k, rid)},
	}
}

func (ln *leafNode[K]) key() K {
	return ln.k
}

func (ln *leafNode[K]) count() int {
	count := 0
	for _, l := range ln.children {
		count += l.count()
	}
	return count
}

func (ln *leafNode[K]) getValue(k K) (RecordID, bool) {
	if ln == nil {
		return ln.zero()
	}

	return ln.getChild(k).getValue()
}

func (ln *leafNode[K]) getChild(k K) *leafEntry[K] {
	for _, l := range ln.children {
		if l.key == k {
			return l
		}
	}

	return nil
}

func (ln *leafNode[K]) insert(f int, k K, rid RecordID) (node[K], bool) {
	c := ln.getChild(k)
	if c != nil {
		// TODO: Create Overflow
		c.append(rid)
		return nil, false
	}

	i := ln.getInsertionIndex(k)
	if i == 0 {
		ln.k = k
	}

	// TODO: Split before inserting
	ln.children = slices.Insert(ln.children, i, newLeafEntry(k, rid))
	if len(ln.children) <= f {
		return nil, false
	}

	return ln.split(), true
}

func (ln *leafNode[K]) getInsertionIndex(k K) int {
	for i := len(ln.children) - 1; i >= 0; i-- {
		j := ln.children[i]
		if k >= j.key {
			return i + 1
		}
	}
	return 0
}

func (ln *leafNode[K]) split() node[K] {
	half := len(ln.children) / 2

	ln2 := &leafNode[K]{
		k:        ln.children[half].key,
		children: ln.children[half:],
		left:     ln,
		right:    ln.right,
	}

	// TODO: Should we copy slices?
	ln.children = slices.Clone(ln.children[:half])
	ln.right = ln2

	return ln2
}

func (ln *leafNode[K]) delete(f int, k K, rid RecordID) {

}

func (ln *leafNode[K]) deleteAll(f int, k K) {
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

func (ln *leafNode[K]) popLargest() *leafEntry[K] {
	largest := ln.children[0]
	ln.children = ln.children[1:]
	return largest
}

func (ln *leafNode[K]) popSmallest() *leafEntry[K] {
	i := len(ln.children) - 1
	smallest := ln.children[i]
	ln.children = ln.children[:i]
	return smallest
}

func (ln *leafNode[K]) getChildForDeletion(k K) (*leafEntry[K], int) {
	for i, l := range ln.children {
		if l.key == k {
			return l, i
		}
	}

	return nil, 0
}

func (ln *leafNode[K]) isMinimum(f int) bool {
	return len(ln.children) <= ln.minimum(f)
}

// TODO: Can we store this?
func (ln *leafNode[K]) minimum(f int) int {
	return ceiling(f, 2) - 1
}

func ceiling(a, b int) int {
	return (a + b - 1) / b
}

func (ln *leafNode[K]) height() int {
	if ln == nil {
		return 0
	}

	return 1
}

func (*leafNode[K]) zero() (RecordID, bool) {
	return RecordID{}, false
}
