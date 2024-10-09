package value

import "errors"

type (
	AggregateErrorMessage string
)

var (
	ErrNotFound                 = errors.New("not found")
	ErrAggregateVersionInvalid  = errors.New("aggregate version invalid")
	ErrAggregateVersionMismatch = errors.New("aggregate version mismatch")
)
