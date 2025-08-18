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

func (bt *BTree[K, V]) FanOut() int {
	return 3
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
		return
	}

	n, ok := bt.root.insert(bt.FanOut(), k, v)
	if !ok {
		return
	}

	bt.newRoot(n)
}

func (bt *BTree[K, V]) newRoot(n node[K, V]) {
	kn := newKeyNode[K, V](bt.root.key())
	kn.children = []node[K, V]{bt.root, n}
	bt.root = kn
}

func (*BTree[K, V]) zero() (V, bool) {
	var v V
	return v, false
}
