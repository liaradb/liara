package btreememory

import "errors"

var (
	ErrEmptyTree = errors.New("empty tree")
	ErrNotFound  = errors.New("not found")
	ErrNoInsert  = errors.New("could not insert")
)
