package btree

import (
	"iter"

	"github.com/liaradb/liaradb/collection/btree/page"
)

type LeafNode struct {
	page page.BTreePage
}

func NewLeafNode(page page.BTreePage) *LeafNode {
	return &LeafNode{
		page: page,
	}
}

func (ln *LeafNode) LeftID() BlockPosition {
	return BlockPosition(ln.page.LowID())
}

func (ln *LeafNode) RightID() BlockPosition {
	return BlockPosition(ln.page.HighID())
}

func (ln *LeafNode) Append(key Key, recordID RecordID) (int16, bool) {
	le := newLeafEntry(key, recordID)
	i, b, ok := ln.page.Append(int16(le.Size()))
	if !ok {
		return 0, false
	}

	le.Write(b)

	return i, true
}

func (ln *LeafNode) Insert(key Key, recordID RecordID) (int16, bool) {
	le := newLeafEntry(key, recordID)
	i := ln.searchIndex(le.key)

	i, b, ok := ln.page.Insert(int16(le.Size()), i)
	if !ok {
		return 0, false
	}

	le.Write(b)

	return i, true
}

func (ln *LeafNode) Child(index int16) (LeafEntry, bool) {
	b, ok := ln.page.Child(index)
	if !ok {
		return LeafEntry{}, false
	}

	le := LeafEntry{}
	le.Read(b)

	return le, true
}

func (ln *LeafNode) Children() iter.Seq[LeafEntry] {
	return func(yield func(LeafEntry) bool) {
		for b := range ln.page.Children() {
			le := LeafEntry{}
			le.Read(b)
			if !yield(le) {
				return
			}
		}
	}
}

// TODO: Handle not found
func (ln *LeafNode) Search(k Key) (RecordID, bool) {
	i := ln.searchIndex(k)

	le, ok := ln.Child(i)
	if !ok {
		return RecordID{}, false
	}

	return le.recordID, true
}

// TODO: Handle not found
func (ln *LeafNode) searchIndex(k Key) int16 {
	var i int16 = 0
	for ke := range ln.Children() {
		if k <= ke.key {
			break
		}

		i++
	}
	return i
}
