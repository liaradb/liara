package btreememory

import "errors"

var (
	ErrAlreadyInitialized = errors.New("already initialized")
	ErrNotFound           = errors.New("not found")
	ErrNoInsert           = errors.New("could not insert")
)
