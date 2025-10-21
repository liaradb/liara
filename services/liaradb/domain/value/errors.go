package value

import (
	"github.com/cardboardrobots/baseerror"
)

var (
	ErrAggregateVersionInvalid = baseerror.ErrInvalidArgument.Wrap("aggregate version invalid")
	// ErrAggregateVersionMismatch = errors.New("aggregate version mismatch")
)
