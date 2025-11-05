package btree

import (
	"cmp"

	"github.com/liaradb/liaradb/storage"
)

type node[K cmp.Ordered, V any] interface {
	key() K
	children() []pair[K, V]
	parentID() storage.BlockID
	rightID() storage.BlockID
	leftID() storage.BlockID
}

type pair[K cmp.Ordered, V any] interface {
	key() K
	value() V
}
