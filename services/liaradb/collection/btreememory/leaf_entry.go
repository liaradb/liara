package btreememory

import "cmp"

type leafEntry[K cmp.Ordered] struct {
	key   K
	value []RecordID
}

func newLeafEntry[K cmp.Ordered](k K, rid RecordID) *leafEntry[K] {
	return &leafEntry[K]{
		key:   k,
		value: []RecordID{rid},
	}
}

func (l *leafEntry[K]) append(rid RecordID) {
	l.value = append(l.value, rid)
}

func (l *leafEntry[K]) getValue() (RecordID, bool) {
	if l == nil {
		return l.zero()
	}

	return l.value[0], true
}

func (l *leafEntry[K]) count() int {
	return len(l.value)
}

func (*leafEntry[K]) zero() (RecordID, bool) {
	return RecordID{}, false
}
