package page

import "errors"

var (
	ErrInvalidCRC = errors.New("invalid CRC")
	ErrNotFound   = errors.New("not found")
	ErrNotPage    = errors.New("not page")
)
