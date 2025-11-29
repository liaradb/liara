package btree

import (
	"github.com/liaradb/liaradb/collection/btree/page"
	"github.com/liaradb/liaradb/storage"
)

type LeafNode struct {
	page page.BTreePage
}

func (ln *LeafNode) LeftID() storage.Offset {
	return ln.page.LowID()
}

func (ln *LeafNode) RightID() storage.Offset {
	return ln.page.HighID()
}

func (ln *LeafNode) Child(index int16) (LeafEntry, bool) {
	_, ok := ln.page.Child(index)
	if !ok {
		return LeafEntry{}, false
	}

	return LeafEntry{}, false
}
