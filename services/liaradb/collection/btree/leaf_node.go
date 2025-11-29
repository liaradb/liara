package btree

import (
	"iter"

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

// TODO: Test this
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

// TODO: Test this
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
