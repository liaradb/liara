package btree

import (
	"cmp"
	"slices"
)

type keyNode[K cmp.Ordered, V any] struct {
	k        K
	level    int
	children []node[K, V]
	left     *keyNode[K, V]
	right    *keyNode[K, V]
}

var _ node[int, int] = (*keyNode[int, int])(nil)

func newKeyNode[K cmp.Ordered, V any](a, b node[K, V]) *keyNode[K, V] {
	return &keyNode[K, V]{
		level:    a.height() + 1,
		children: []node[K, V]{a, b},
	}
}

func (kn *keyNode[K, V]) key() K {
	return kn.k
}

func (kn *keyNode[K, V]) count() int {
	count := 0
	for _, l := range kn.children {
		count += l.count()
	}
	return count
}

func (kn *keyNode[K, V]) getValue(k K) (V, bool) {
	if kn == nil || kn.count() == 0 {
		return kn.zero()
	}

	return kn.getChild(k).getValue(k)
}

func (kn *keyNode[K, V]) getChild(k K) node[K, V] {
	a := kn.children[0]

	l := len(kn.children)
	for i := 1; i < l; i++ {
		b := kn.children[i]
		if k < b.key() {
			return a
		}

		a = b
	}

	return a
}

func (kn *keyNode[K, V]) insert(f int, k K, v V) (node[K, V], bool) {
	n, ok := kn.getChild(k).insert(f, k, v)
	if !ok {
		return nil, false
	}

	return kn.insertNode(f, k, n)
}

func (kn *keyNode[K, V]) insertNode(f int, k K, n node[K, V]) (node[K, V], bool) {
	i := kn.getInsertionIndex(n.key())
	if i == 0 {
		kn.k = k
	}

	// TODO: Split before inserting
	kn.children = slices.Insert(kn.children, i, n)
	if len(kn.children) <= f {
		return nil, false
	}

	return kn.split(), true
}

func (kn *keyNode[K, V]) getInsertionIndex(k K) int {
	for i := len(kn.children) - 1; i >= 0; i-- {
		j := kn.children[i]
		if k >= j.key() {
			return i + 1
		}
	}
	return 0
}

func (kn *keyNode[K, V]) split() node[K, V] {
	half := len(kn.children) / 2

	kn2 := &keyNode[K, V]{
		k:        kn.children[half].key(),
		children: kn.children[half:],
		left:     kn,
		right:    kn.right,
	}

	// TODO: Should we copy slices?
	kn.children = slices.Clone(kn.children[:half])
	kn.right = kn2

	return kn2
}

func (kn *keyNode[K, V]) delete(f int, k K, v V) {

}

func (kn *keyNode[K, V]) deleteAll(f int, k K) {

}

func (kn *keyNode[K, V]) height() int {
	if kn == nil {
		return 0
	}

	return kn.level
}

func (*keyNode[K, V]) zero() (V, bool) {
	var v V
	return v, false
}
