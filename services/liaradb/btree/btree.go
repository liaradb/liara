package btree

type BTree[K comparable, V any] struct {
	root node[K, V]
}

func (bt *BTree[K, V]) Height() int {
	if bt.root == nil {
		return 0
	}
	return bt.root.height()
}

func (bt *BTree[K, V]) FanOut() int {
	return 2
}

func (bt *BTree[K, V]) getValue(k K) (V, bool) {
	if bt.root == nil {
		return bt.zero()
	}
	return bt.root.getValue(k)
}

func (bt *BTree[K, V]) insert(k K, v V) {
	if bt.root == nil {
		bt.root = newLeafNode(k, v)
	} else {
		bt.root.insert(k, v)
	}
}

func (*BTree[K, V]) zero() (V, bool) {
	var v V
	return v, false
}
