package replay

import "errors"

var (
	ErrActionUnknown     = errors.New("action unknown")
	ErrCollectionUnknown = errors.New("collection unknown")
)
