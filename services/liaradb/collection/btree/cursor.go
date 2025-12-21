package btree

import (
	"github.com/liaradb/liaradb/collection/btree/value"
	"github.com/liaradb/liaradb/storage"
)

type Cursor struct {
	insert
	level
	search
}

type (
	Key            = value.Key
	RecordID       = value.RecordID
	RecordPosition = value.RecordPosition
)

func NewCursor(s *storage.Storage) *Cursor {
	return &Cursor{
		insert: newInsert(s),
		level:  newLevel(s),
		search: newSearch(s),
	}
}
