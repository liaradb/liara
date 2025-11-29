package btree

import (
	"iter"

	"github.com/liaradb/liaradb/collection/btree/page"
)

type KeyNode struct {
	page page.BTreePage
}

// TODO: Test this
// TODO: Change to bool instead of error
func (kn *KeyNode) Children() iter.Seq2[KeyEntry, error] {
	return func(yield func(KeyEntry, error) bool) {
		for b := range kn.page.Children() {
			ke := KeyEntry{}
			if err := ke.Read(b); err != nil {
				yield(KeyEntry{}, err)
				return
			}

			if !yield(ke, nil) {
				return
			}
		}
	}
}

// TODO: Test this
// TODO: Change to bool instead of error
func (kn *KeyNode) Search(k Key) (BlockPosition, error) {
	p := BlockPosition(kn.page.LowID())
	for ke, err := range kn.Children() {
		if err != nil {
			return 0, err
		}

		if k < ke.key {
			break
		}

		p = ke.block
	}
	return p, nil
}
