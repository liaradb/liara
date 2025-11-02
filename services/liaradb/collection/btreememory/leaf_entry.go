package btreememory

import "cmp"

type leafEntry[K cmp.Ordered, V any] struct {
	key   K
	value []V
}

func newLeafEntry[K cmp.Ordered, V any](k K, v V) *leafEntry[K, V] {
	return &leafEntry[K, V]{
		key:   k,
		value: []V{v},
	}
}

func (l *leafEntry[K, V]) append(v V) {
	l.value = append(l.value, v)
}

func (l *leafEntry[K, V]) getValue() (V, bool) {
	if l == nil {
		return l.zero()
	}

	return l.value[0], true
}

func (l *leafEntry[K, V]) count() int {
	return len(l.value)
}

func (*leafEntry[K, V]) zero() (V, bool) {
	var v V
	return v, false
}
