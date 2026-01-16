package btree

import (
	"errors"

	"github.com/cardboardrobots/baseerror"
)

var (
	ErrExists        = baseerror.ErrAlreadyExists
	ErrLevelMismatch = errors.New("level mismatch")
	ErrNotFound      = baseerror.ErrNotFound
	ErrNoInsert      = errors.New("could not insert")
	ErrNoUpdate      = errors.New("could not update")
	ErrTypeMismatch  = errors.New("type mismatch")
)
