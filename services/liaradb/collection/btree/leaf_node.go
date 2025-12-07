package btree

import (
	"iter"

	"github.com/liaradb/liaradb/collection/btree/page"
	"github.com/liaradb/liaradb/encoder/raw"
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

	// TODO: Change to bool instead of error
	if err := le.Write(raw.NewBufferFromSlice(b)); err != nil {
		return 0, false
	}

	return i, true
}

func (ln *LeafNode) Insert(key Key, recordID RecordID) (int16, bool) {
	le := newLeafEntry(key, recordID)
	i, err := ln.searchIndex(le.key)
	if err != nil {
		return 0, false
	}

	i, b, ok := ln.page.Insert(int16(le.Size()), i)
	if !ok {
		return 0, false
	}

	// TODO: Change to bool instead of error
	if err := le.Write(raw.NewBufferFromSlice(b)); err != nil {
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
	if err := le.Read(raw.NewBufferFromSlice(b)); err != nil {
		return LeafEntry{}, err
	}

	return le, nil
}

// TODO: Change to bool instead of error
func (ln *LeafNode) Children() iter.Seq2[LeafEntry, error] {
	return func(yield func(LeafEntry, error) bool) {
		for b := range ln.page.Children() {
			le := LeafEntry{}
			if err := le.Read(raw.NewBufferFromSlice(b)); err != nil {
				yield(LeafEntry{}, err)
				return
			}

			if !yield(le, nil) {
				return
			}
		}
	}
}

func (ln *LeafNode) Search(k Key) (RecordID, error) {
	i, err := ln.searchIndex(k)
	if err != nil {
		return RecordID{}, err
	}

	le, err := ln.Child(i)
	if err != nil {
		return RecordID{}, err
	}

	return le.recordID, nil
}

func (ln *LeafNode) searchIndex(k Key) (int16, error) {
	var i int16 = 0
	for ke, err := range ln.Children() {
		if err != nil {
			return 0, err
		}

		if k <= ke.key {
			break
		}

		i++
	}
	return i, nil
}
