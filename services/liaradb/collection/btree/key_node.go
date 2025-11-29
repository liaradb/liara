package btree

import (
	"iter"

	"github.com/liaradb/liaradb/collection/btree/page"
)

type KeyNode struct {
	page page.BTreePage
}

func (kn *KeyNode) Keys() iter.Seq[string] {
	return func(yield func(string) bool) {
	}
}
