package storage

import "errors"

var (
	ErrNotInitialized = errors.New("not initialized")
	ErrRequestClosed  = errors.New("request closed")
)
