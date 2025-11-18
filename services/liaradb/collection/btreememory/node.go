package btreememory

import "cmp"

type node[K cmp.Ordered] interface {
	key() K
	getValue(k K) (RecordID, bool)
	insert(fanout int, k K, rid RecordID) (node[K], bool)
	delete(fanout int, k K, rid RecordID)
	deleteAll(fanout int, k K)
	height() int
	count() int
}
