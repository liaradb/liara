package btree

import "errors"

var (
	ErrNotFound = errors.New("not found")
	ErrNoInsert = errors.New("could not insert")
)
