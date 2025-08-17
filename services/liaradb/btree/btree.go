package btree

type BTree[K comparable, V any] struct {
	root *keyNode[K, V]
}

func (bt *BTree[K, V]) Height() int {
	return bt.root.height()
}

func (bt *BTree[K, V]) FanOut() int {
	return 3
}

func (bt *BTree[K, V]) getValue(k K) (V, bool) {
	return bt.root.getValue(k)
}

func (bt *BTree[K, V]) insert(k K, v V) {
	if bt.root == nil {
		bt.root = &keyNode[K, V]{}
	}
	bt.root.insert(k, v)
}
