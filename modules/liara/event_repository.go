package liara

import (
	"context"
	"iter"
	"time"
)

type EventRepository interface {
	Get(context.Context, AggregateID) iter.Seq2[Event, error]
	GetAfterGlobalVersion(context.Context, GlobalVersion, Limit) iter.Seq2[Event, error]
	GetByAggregateIDAndName(context.Context, AggregateID, AggregateName) iter.Seq2[Event, error]
	Append(context.Context, ...AppendEvent) error
}

type AppendEvent struct {
	ID            EventID       // The ID of the Event, used for de-duplication
	AggregateID   AggregateID   // The ID of the Aggregate to which this Event applies
	Version       Version       // The Version of the Aggregate
	PartitionID   PartitionID   // The ID to partition Events
	UserID        UserID        // The ID of the User issuing the Command
	CorrelationID CorrelationID // The ID of the entire Command and Event chain
	Time          time.Time     // The Time this Event was created
	AggregateName AggregateName // The Name of the Aggregate
	Name          EventName     // The Name of the Event
	Schema        Schema        // The Schema for the internal data
	Data          []byte        // The internal data of the Event
}

func (ae *AppendEvent) Valid() error {
	if ae.Version < 1 {
		return ErrAggregateVersionInvalid
	}

	return nil
}
