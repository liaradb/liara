package btree

type leafEntry[K comparable, V any] struct {
	key   K
	value V
}

func newLeafEntry[K comparable, V any](k K, v V) *leafEntry[K, V] {
	return &leafEntry[K, V]{
		key:   k,
		value: v,
	}
}

func (l *leafEntry[K, V]) getValue() (V, bool) {
	if l == nil {
		var v V
		return v, false
	}

	return l.value, true
}
