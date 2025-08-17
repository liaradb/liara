package btree

type node[K comparable, V any] interface {
	key() K
	getValue(k K) (V, bool)
	insert(k K, v V)
	height() int
}
