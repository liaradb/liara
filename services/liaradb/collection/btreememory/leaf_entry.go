package btreememory

import "cmp"

type leafEntry[K cmp.Ordered] struct {
	key   K
	value RecordID
}

func newLeafEntry[K cmp.Ordered](k K, rid RecordID) *leafEntry[K] {
	return &leafEntry[K]{
		key:   k,
		value: rid,
	}
}

func (l *leafEntry[K]) getValue() (RecordID, bool) {
	if l == nil {
		return l.zero()
	}

	return l.value, true
}

func (*leafEntry[K]) zero() (RecordID, bool) {
	return RecordID{}, false
}
