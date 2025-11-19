package btreememory

import (
	"cmp"
	"slices"

	"github.com/liaradb/liaradb/storage"
)

type leafNode[K cmp.Ordered] struct {
	storage  Storage[K]
	i        storage.Offset
	k        K
	children []*leafEntry[K]
	leftID   storage.Offset
	rightID  storage.Offset
}

var _ node[int] = (*leafNode[int])(nil)

func newLeafNode[K cmp.Ordered](s Storage[K], k K, rid RecordID) *leafNode[K] {
	return &leafNode[K]{
		i:        nextID(),
		storage:  s,
		k:        k,
		children: []*leafEntry[K]{newLeafEntry(k, rid)},
	}
}

func (ln *leafNode[K]) key() K             { return ln.k }
func (ln *leafNode[K]) id() storage.Offset { return ln.i }
func (ln *leafNode[K]) count() int         { return len(ln.children) }

func (ln *leafNode[K]) getChild(k K) (storage.Offset, bool) {
	return 0, false
}

func (ln *leafNode[K]) getValue(k K) (RecordID, bool) {
	if ln == nil {
		return RecordID{}, false
	}

	c, ok := ln.getChildValue(k)
	if !ok {
		return RecordID{}, false
	}

	return c.value, true
}

func (ln *leafNode[K]) getChildValue(k K) (*leafEntry[K], bool) {
	for _, l := range ln.children {
		if l.key == k {
			return l, true
		}
	}

	return nil, false
}

func (ln *leafNode[K]) insert(f int, k K, rid RecordID) (*leafNode[K], bool) {
	c, ok := ln.getChildValue(k)
	if ok {
		// TODO: Create Overflow
		c.value = rid
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

func (ln *leafNode[K]) split() *leafNode[K] {
	half := len(ln.children) / 2

	ln2 := &leafNode[K]{
		i:        nextID(),
		k:        ln.children[half].key,
		children: ln.children[half:],
		leftID:   ln.i,
		rightID:  ln.rightID,
	}

	// TODO: Should we copy slices?
	ln.children = slices.Clone(ln.children[:half])
	ln.rightID = ln2.i

	return ln2
}

func (ln *leafNode[K]) delete(f int, k K, rid RecordID) {

}

func (ln *leafNode[K]) deleteAll(f int, k K) {
	c, i := ln.getChildForDeletion(k)
	if c == nil {
		return
	}

	// if ln.isMinimum(f) {
	// 	// TODO: Rebalance
	// 	if ln.left != nil && !ln.left.isMinimum(f) {
	// 		// Borrow Left
	// 		e := ln.left.popSmallest()
	// 		// TODO: How do we handle overflow?
	// 		ln.insert(f, e.key, e.value[0])
	// 		// Pull smallest from Left
	// 		// Update This Key
	// 		// Key change propagates
	// 	} else if ln.right != nil && !ln.right.isMinimum(f) {
	// 		// Borrow Right
	// 		e := ln.right.popLargest()
	// 		// TODO: How do we handle overflow?
	// 		ln.insert(f, e.key, e.value[0])
	// 		// Pull largest from Right
	// 		// Update Right Key
	// 		// Key changes propagates
	// 	} else if ln.left != nil {
	// 		// Merge Left
	// 		// Move children to Left
	// 		// Delete node
	// 		// Deletion propagates
	// 	} else if ln.right != nil {
	// 		// Merge Right
	// 		// Move children to Right
	// 		// Update Right Key
	// 		// Key change propagates
	// 		// Delete node
	// 		// Deletion propagates
	// 	} else {
	// 		// Delete
	// 	}
	// } else {
	// 	// Delete
	// }
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
