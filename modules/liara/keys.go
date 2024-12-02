package liara

import "github.com/google/uuid"

type (
	AggregateID   string
	CorrelationID string
	EventID       string
	OutboxID      string
	RequestID     string
	PartitionID   int32
	UserID        string
)

func (a AggregateID) String() string   { return string(a) }
func (c CorrelationID) String() string { return string(c) }
func (e EventID) String() string       { return string(e) }
func (o OutboxID) String() string      { return string(o) }
func (r RequestID) String() string     { return string(r) }
func (p PartitionID) Value() int32     { return int32(p) }
func (u UserID) String() string        { return string(u) }

func NewAggregateID() AggregateID {
	return AggregateID(uuid.NewString())
}

func NewEventID() EventID {
	return EventID(uuid.NewString())
}
