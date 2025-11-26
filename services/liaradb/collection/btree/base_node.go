package btree

import "github.com/liaradb/liaradb/storage"

type BaseNode[K Key, V any] struct {
	b   *storage.Buffer
	k   K
	pID storage.BlockID
	rID storage.BlockID
	lID storage.BlockID
}

func (bn *BaseNode[K, V]) key() K {
	return bn.k
}

func (bn *BaseNode[K, V]) children() []pair[K, V] {
	return nil
}

func (bn *BaseNode[K, V]) parentID() storage.BlockID {
	return bn.pID
}

func (bn *BaseNode[K, V]) rightID() storage.BlockID {
	return bn.rID
}

func (bn *BaseNode[K, V]) leftID() storage.BlockID {
	return bn.lID
}
