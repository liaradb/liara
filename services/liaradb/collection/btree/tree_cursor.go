package btree

import "cmp"

type TreeCursor[K cmp.Ordered, V any] struct {
}

func (tc *TreeCursor[K, V]) Current() (node[K, V], bool) {
	return nil, false
}

func (tc *TreeCursor[K, V]) Left() (node[K, V], bool) {
	return nil, false
}

func (tc *TreeCursor[K, V]) Parent() (node[K, V], bool) {
	return nil, false
}

func (tc *TreeCursor[K, V]) Right() (node[K, V], bool) {
	return nil, false
}
