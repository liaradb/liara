package page

import "errors"

var (
	ErrInsufficientSpace = errors.New("insufficient space")
	ErrInvalidCRC        = errors.New("invalid CRC")
	ErrNotPage           = errors.New("not page")
)
