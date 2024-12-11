package liara

import (
	"context"
	"iter"
)

type EventRepository interface {
	Get(context.Context, TenantID, AggregateID) iter.Seq2[Event, error]
	GetAfterGlobalVersion(context.Context, TenantID, GlobalVersion, []PartitionID, Limit) iter.Seq2[Event, error]
	GetByAggregateIDAndName(context.Context, TenantID, AggregateID, AggregateName) iter.Seq2[Event, error]
	GetByOutbox(context.Context, TenantID, OutboxID, Limit) iter.Seq2[Event, error]
	Append(context.Context, TenantID, AppendOptions, ...AppendEvent) error
	TestIdempotency(context.Context, TenantID, RequestID) (bool, error)
}

type AppendEvent struct {
	ID            EventID       // The ID of the Event, used for de-duplication
	AggregateID   AggregateID   // The ID of the Aggregate to which this Event applies
	Version       Version       // The Version of the Aggregate
	PartitionID   PartitionID   // The ID to partition Events
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
