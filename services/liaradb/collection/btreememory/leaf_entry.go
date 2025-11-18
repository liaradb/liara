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
