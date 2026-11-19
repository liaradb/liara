package storage

import "errors"

var (
	ErrNotInitialized = errors.New("not initialized")
	ErrNotDirty       = errors.New("not dirty")
)
