package btree

import "github.com/liaradb/liaradb/storage"

type Cursor struct {
	insert
	level
	search
}

func NewCursor(s *storage.Storage) *Cursor {
	return &Cursor{
		insert: newInsert(s),
		level:  newLevel(s),
		search: newSearch(s),
	}
}
