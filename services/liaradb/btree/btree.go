package btree

import "cmp"

type BTree[K cmp.Ordered, V any] struct {
	root node[K, V]
}

func (bt *BTree[K, V]) Height() int {
	if bt.root == nil {
		return 0
	}
	return bt.root.height()
}

func (bt *BTree[K, V]) Count() int {
	if bt.root == nil {
		return 0
	}
	return bt.root.count()
}

func (bt *BTree[K, V]) FanOut() int {
	return 3
}

func (bt *BTree[K, V]) GetValue(k K) (V, bool) {
	if bt.root == nil {
		return bt.zero()
	}
	return bt.root.getValue(k)
}

func (bt *BTree[K, V]) Insert(k K, v V) {
	if bt.root == nil {
		bt.root = newLeafNode(k, v)
		return
	}

	n, ok := bt.root.insert(bt.FanOut(), k, v)
	if !ok {
		return
	}

	bt.root = newKeyNode(bt.root, n)
}

func (bt *BTree[K, V]) DeleteAll(k K) {
	if bt.root == nil {
		return
	}

	bt.root.deleteAll(k)
}

func (bt *BTree[K, V]) DeleteValue(k K, v V) {

}

func (*BTree[K, V]) zero() (V, bool) {
	var v V
	return v, false
}
