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

func (ln *LeafNode) Append(le LeafEntry) (int16, bool) {
	i, b, ok := ln.page.Append(int16(le.Size()))
	if !ok {
		return 0, false
	}

	// TODO: Change to bool instead of error
	if err := le.Write(b); err != nil {
		return 0, false
	}

	return i, true
}

// TODO: Change to bool instead of error
func (ln *LeafNode) Child(index int16) (LeafEntry, error) {
	b, ok := ln.page.Child(index)
	if !ok {
		return LeafEntry{}, ErrNotFound
	}

	le := LeafEntry{}
	if err := le.Read(b); err != nil {
		return LeafEntry{}, err
	}

	return le, nil
}

// TODO: Change to bool instead of error
func (ln *LeafNode) Children() iter.Seq2[LeafEntry, error] {
	return func(yield func(LeafEntry, error) bool) {
		for b := range ln.page.Children() {
			le := LeafEntry{}
			if err := le.Read(b); err != nil {
				yield(LeafEntry{}, err)
				return
			}

			if !yield(le, nil) {
				return
			}
		}
	}
}
