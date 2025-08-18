package btree

type leafNode[K comparable, V any] struct {
	k        K
	children []*leafEntry[K, V]
}

var _ node[int, int] = (*leafNode[int, int])(nil)

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
		return ln.zero()
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
	ln.children = append(ln.children, newLeafEntry(k, v))
}

// TODO: Should we copy slices?
func (ln *leafNode[K, V]) split() *leafNode[K, V] {
	half := len(ln.children) / 2

	ln2 := &leafNode[K, V]{
		k:        ln.children[half].key,
		children: ln.children[half:],
	}

	ln.children = ln.children[:half]

	return ln2
}

func (ln *leafNode[K, V]) height() int {
	if ln == nil {
		return 0
	}

	return 1
}

func (*leafNode[K, V]) zero() (V, bool) {
	var v V
	return v, false
}
