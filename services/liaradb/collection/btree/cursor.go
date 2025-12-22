package btree

import (
	"github.com/liaradb/liaradb/collection/btree/value"
	"github.com/liaradb/liaradb/storage"
	"github.com/liaradb/liaradb/storage/link"
)

type Cursor struct {
	insert
	level
	search
}

type (
	Key            = value.Key
	RecordID       = link.RecordLocator
	RecordPosition = link.RecordPosition
)

func NewCursor(s *storage.Storage) *Cursor {
	return &Cursor{
		insert: newInsert(s),
		level:  newLevel(s),
		search: newSearch(s),
	}
}
