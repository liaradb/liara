package btree

import "github.com/liaradb/liaradb/collection/log"

type Cursor struct {
	insert
	level
	search
	maximum
}

func NewCursor(l log.Log) *Cursor {
	return &Cursor{
		insert:  newInsert(l),
		level:   newLevel(l),
		search:  newSearch(l),
		maximum: newMaximum(l),
	}
}
