package btree

import "github.com/liaradb/liaradb/encoder/page"

type LeafNode[K Key, V page.Serializer] struct {
	Header LeafNodeHeader[K]
	Keys   []KeyEntry[K]
	Values []V
}
