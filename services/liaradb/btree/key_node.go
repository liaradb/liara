package btree

type keyNode[K comparable, V any] struct {
	k        K
	level    int
	children []node[K, V]
}

var _ node[int, int] = (*keyNode[int, int])(nil)

func newKeyNode[K comparable, V any](k K) *keyNode[K, V] {
	return &keyNode[K, V]{
		k: k,
	}
}

func (kn *keyNode[K, V]) key() K {
	return kn.k
}

func (kn *keyNode[K, V]) count() int {
	return len(kn.children)
}

func (kn *keyNode[K, V]) getValue(k K) (V, bool) {
	if kn == nil {
		return kn.zero()
	}

	var child node[K, V]
	for _, ln := range kn.children {
		if ln.key() == k {
			child = ln
			break
		}
	}
	if child == nil {
		return kn.zero()
	}

	return child.getValue(k)
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
