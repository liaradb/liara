package btreememory

import (
	"cmp"
	"fmt"
)

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

func (le leafEntry[K]) String() string {
	return fmt.Sprintf("(%v -> X)", le.key)
}
