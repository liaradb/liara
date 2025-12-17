package btree

import "github.com/liaradb/liaradb/storage"

type Cursor struct {
	insert
	search
}

func NewCursor(s *storage.Storage) *Cursor {
	return &Cursor{
		insert: newInsert(s),
		search: newSearch(s),
	}
}
