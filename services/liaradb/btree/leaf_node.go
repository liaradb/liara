package btree

type leafNode[K comparable, V any] struct {
	k        K
	children []*leafEntry[K, V]
}

func newLeafNode[K comparable, V any](k K, v V) *leafNode[K, V] {
	return &leafNode[K, V]{
		k:        k,
		children: []*leafEntry[K, V]{newLeafEntry(k, v)},
	}
}

func (ln *leafNode[K, V]) key() K {
	return ln.k
}

func (ln *leafNode[K, V]) getValue(k K) (V, bool) {
	if ln == nil {
		var v V
		return v, false
	}

	var child *leafEntry[K, V]
	for _, l := range ln.children {
		if l.key == k {
			child = l
			break
		}
	}
	return child.getValue()
}

func (ln *leafNode[K, V]) insert(k K, v V) {
}

func (ln *leafNode[K, V]) height() int {
	if ln == nil {
		return 0
	}

	return 1
}
