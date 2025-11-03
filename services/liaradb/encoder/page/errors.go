package page

import "errors"

var (
	ErrInvalidCRC = errors.New("invalid CRC")
	ErrNotPage    = errors.New("not page")
)
