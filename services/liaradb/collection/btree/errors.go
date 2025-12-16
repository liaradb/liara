package btree

import "errors"

var (
	ErrLevelMismatch = errors.New("level mismatch")
	ErrNotFound      = errors.New("not found")
	ErrNoInsert      = errors.New("could not insert")
	ErrTypeMismatch  = errors.New("type mismatch")
)
