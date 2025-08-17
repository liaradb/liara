package btree

type keyNode[K comparable, V any] struct {
	level    int
	children []node[K, V]
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
	return child.getValue(k)
}

func (kn *keyNode[K, V]) insert(k K, v V) {
	kn.children = append(kn.children, newLeafNode(k, v))
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
