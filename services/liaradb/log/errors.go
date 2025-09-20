package log

import "errors"

var (
	ErrInvalidCRC        = errors.New("invalid CRC")
	ErrNotPage           = errors.New("not page")
	ErrInsufficientSpace = errors.New("insufficient space")
)
