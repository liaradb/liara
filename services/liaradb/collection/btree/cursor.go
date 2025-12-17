package btree

import "github.com/liaradb/liaradb/storage"

type Cursor struct {
	insertCursor
	searchCursor
}

func NewCursor(s *storage.Storage) *Cursor {
	return &Cursor{
		insertCursor: newInsertCursor(s),
		searchCursor: newSearchCursor(s),
	}
}
