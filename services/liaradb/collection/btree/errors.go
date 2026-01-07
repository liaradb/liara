package btree

import "errors"

var (
	ErrExists        = errors.New("already exists at key")
	ErrLevelMismatch = errors.New("level mismatch")
	ErrNotFound      = errors.New("not found")
	ErrNoInsert      = errors.New("could not insert")
	ErrTypeMismatch  = errors.New("type mismatch")
)
