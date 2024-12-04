package service

import (
	"context"
	"iter"

	"github.com/cardboardrobots/liara_service/feature/eventsource/domain/entity"
	"github.com/cardboardrobots/liara_service/feature/eventsource/domain/value"
)

type EventRepository interface {
	CreateTable(context.Context, value.TenantID) error
	CreateIndex(context.Context, value.TenantID) error
	DropTable(context.Context, value.TenantID) error
	Get(context.Context, value.TenantID, value.AggregateID) iter.Seq2[entity.Event, error]
	GetAfterGlobalVersion(context.Context, value.TenantID, value.GlobalVersion, value.PartitionRange, value.Limit) iter.Seq2[entity.Event, error]
	GetByAggregateIDAndName(context.Context, value.TenantID, value.AggregateID, value.AggregateName) iter.Seq2[entity.Event, error]
	Append(context.Context, value.TenantID, AppendEvent) error
}

type AppendEvent struct {
	ID            value.EventID        // The ID of the Event, used for de-duplication
	AggregateName value.AggregateName  // The Name of the Aggregate
	AggregateID   value.AggregateID    // The ID of the Aggregate to which this Event applies
	Version       value.Version        // The Version of the Aggregate
	PartitionID   value.PartitionID    // The ID to partition Events
	Name          value.EventName      // The Name of the Event
	Schema        value.Schema         // The Schema for the internal data
	Metadata      entity.EventMetadata // The Metadata of the Event
	Data          []byte               // The internal data of the Event
}

func (ae *AppendEvent) Valid() error {
	if ae.Version < 1 {
		return value.ErrAggregateVersionInvalid
	}

	return nil
}
