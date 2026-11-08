package value

import "github.com/cardboardrobots/baseerror"

var (
	ErrNoEvents                = baseerror.ErrInvalidArgument.Wrap("no events")
	ErrRequestIDExists         = baseerror.ErrInvalidArgument.Wrap("request id exists")
	ErrAggregateVersionInvalid = baseerror.ErrInvalidArgument.Wrap("aggregate version invalid")
	// ErrAggregateVersionMismatch = errors.New("aggregate version mismatch")
)
