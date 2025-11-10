package storage

import "errors"

var (
	ErrInvalidRequest = errors.New("invalid request")
	ErrNotInitialized = errors.New("not initialized")
	ErrNotDirty       = errors.New("not dirty")
)
