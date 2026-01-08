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
	ErrTypeMismatch  = errors.New("type mismatch")
)
