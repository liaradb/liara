package btree

import "cmp"

type leafEntry[K cmp.Ordered, V any] struct {
	key   K
	value V
}

func newLeafEntry[K cmp.Ordered, V any](k K, v V) *leafEntry[K, V] {
	return &leafEntry[K, V]{
		key:   k,
		value: v,
	}
}

func (l *leafEntry[K, V]) getValue() (V, bool) {
	if l == nil {
		return l.zero()
	}

	return l.value, true
}

func (*leafEntry[K, V]) zero() (V, bool) {
	var v V
	return v, false
}
