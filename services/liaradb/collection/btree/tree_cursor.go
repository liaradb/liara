package btree

import (
	"cmp"

	"github.com/liaradb/liaradb/storage"
)

type TreeCursor[K cmp.Ordered, V any] struct {
	current node[K, V]
}

func (tc *TreeCursor[K, V]) GetNode(id storage.BlockID) (node[K, V], error) {
	return nil, nil
}

func (tc *TreeCursor[K, V]) Current() (node[K, V], error) {
	return tc.current, nil
}

func (tc *TreeCursor[K, V]) Left() (node[K, V], error) {
	return tc.GetNode(tc.current.leftID())
}

func (tc *TreeCursor[K, V]) Parent() (node[K, V], error) {
	return tc.GetNode(tc.current.parentID())
}

func (tc *TreeCursor[K, V]) Right() (node[K, V], error) {
	return tc.GetNode(tc.current.rightID())
}
