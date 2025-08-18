package btree

type node[K comparable, V any] interface {
	key() K
	getValue(k K) (V, bool)
	insert(fanout int, k K, v V) (node[K, V], bool)
	height() int
	count() int
}
