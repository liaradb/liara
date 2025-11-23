package btreememory

import (
	"cmp"

	"github.com/liaradb/liaradb/storage"
)

type node[K cmp.Ordered] interface {
	id() storage.Offset
	firstKey() K
	isKeyNode() bool
	isLeafNode() bool
	getChild(k K) (storage.Offset, bool)
	getValue(k K) (RecordID, bool)
	// insert(fanout int, k K, rid RecordID) (node[K], bool)
	delete(fanout int, k K, rid RecordID)
	deleteAll(fanout int, k K)
	height() int
	count() int
	String() string
}
