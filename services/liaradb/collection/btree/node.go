package btree

import (
	"github.com/liaradb/liaradb/storage"
)

type node[K Key, V any] interface {
	key() K
	children() []pair[K, V]
	parentID() storage.BlockID
	rightID() storage.BlockID
	leftID() storage.BlockID
}

type pair[K Key, V any] interface {
	key() K
	value() V
}
