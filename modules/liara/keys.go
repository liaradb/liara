package liara

import "github.com/google/uuid"

type (
	AggregateID   string
	CorrelationID string
	EventID       string
	IdempotenceID string
	OutboxID      string
	PartitionID   string
	UserID        string
)

func (a AggregateID) String() string   { return string(a) }
func (c CorrelationID) String() string { return string(c) }
func (e EventID) String() string       { return string(e) }
func (i IdempotenceID) String() string { return string(i) }
func (o OutboxID) String() string      { return string(o) }
func (p PartitionID) String() string   { return string(p) }
func (u UserID) String() string        { return string(u) }

func NewAggregateID() AggregateID {
	return AggregateID(uuid.NewString())
}

func NewEventID() EventID {
	return EventID(uuid.NewString())
}
