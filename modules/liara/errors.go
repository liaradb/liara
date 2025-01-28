package liara

import (
	"github.com/cardboardrobots/baseerror"
)

var (
	ErrNotFound                 = baseerror.ErrNotFound
	ErrAggregateVersionInvalid  = baseerror.ErrInvalidArgument.Wrap("aggregate version invalid")
	ErrAggregateVersionMismatch = baseerror.ErrInvalidArgument.Wrap("aggregate version mismatch")
)
