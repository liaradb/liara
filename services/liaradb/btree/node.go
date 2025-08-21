package btree

import "cmp"

type node[K cmp.Ordered, V any] interface {
	key() K
	getValue(k K) (V, bool)
	insert(fanout int, k K, v V) (node[K, V], bool)
	delete(fanout int, k K, v V)
	deleteAll(fanout int, k K)
	height() int
	count() int
	setParent(n *keyNode[K, V])
}
