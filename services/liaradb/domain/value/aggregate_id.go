package value

import "github.com/google/uuid"

type AggregateID string

func (a AggregateID) String() string { return string(a) }

func NewAggregateID() AggregateID {
	return AggregateID(uuid.NewString())
}
