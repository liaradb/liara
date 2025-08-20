package btree

import "cmp"

type keyNode[K cmp.Ordered, V any] struct {
	k        K
	level    int
	children []node[K, V]
}

var _ node[int, int] = (*keyNode[int, int])(nil)

func newKeyNode[K cmp.Ordered, V any](a, b node[K, V]) *keyNode[K, V] {
	return &keyNode[K, V]{
		children: []node[K, V]{a, b},
	}
}

func (kn *keyNode[K, V]) key() K {
	return kn.k
}

func (kn *keyNode[K, V]) count() int {
	return len(kn.children)
}

func (kn *keyNode[K, V]) getValue(k K) (V, bool) {
	if kn == nil || kn.count() == 0 {
		return kn.zero()
	}

	return kn.getChild(k).getValue(k)
}

func (kn *keyNode[K, V]) getChild(k K) node[K, V] {
	child := kn.children[0]
	for i := range len(kn.children) - 1 {
		c := kn.children[i+1]
		if k >= c.key() {
			child = c
		} else {
			break
		}
	}
	return child
}

func (kn *keyNode[K, V]) insert(f int, k K, v V) (node[K, V], bool) {
	kn.children = append(kn.children, newLeafNode(k, v))
	return nil, false
}

func (kn *keyNode[K, V]) height() int {
	if kn == nil {
		return 0
	}

	return kn.level + 1
}

func (*keyNode[K, V]) zero() (V, bool) {
	var v V
	return v, false
}
