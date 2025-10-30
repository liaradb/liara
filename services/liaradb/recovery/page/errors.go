package page

import "errors"

var (
	ErrInsufficientSpace = errors.New("insufficient space")
	ErrNotPage           = errors.New("not page")
)
