package liara

import (
	"context"
	"encoding/json"
	"iter"
	"time"

	"github.com/cardboardrobots/eventsource/entity"
	"github.com/cardboardrobots/eventsource/value"
)

type EventRepository interface {
	Get(context.Context, value.AggregateID) iter.Seq2[entity.Event, error]
	GetAfterGlobalVersion(context.Context, value.GlobalVersion, value.Limit) iter.Seq2[entity.Event, error]
	GetByAggregateIDAndName(context.Context, value.AggregateID, value.AggregateName) iter.Seq2[entity.Event, error]
	Append(context.Context, ...AppendEvent) error
}

type (
	AppendEvent struct {
		AggregateName value.AggregateName // The name of the Aggregate
		ID            value.EventID       // The unique ID of the Event
		AggregateID   value.AggregateID   // The ID of the Aggregate to which this Event applies
		Version       value.Version       // The Version of the Aggregate
		Name          value.EventName     // The unique GetName of the Event
		CorrelationID value.CorrelationID // The ID of the entire Command and Event chain
		IdempotenceID value.IdempotenceID // The ID to de-duplicate Events
		PartitionID   value.PartitionID   // The ID to partition Events
		UserID        value.UserID        // The ID of the User issuing the Command
		Time          time.Time           // The Time this Event was created
		Schema        value.Schema        // The schema for the internal data
		Data          []byte              // The internal data of the Event
	}

	EventData interface {
		EventName() string
		AggregateName() string
		Schema() string
	}

	EventOptions[U ~string] struct {
		EventID       value.EventID
		Time          time.Time
		AggregateID   U
		Version       value.Version
		CorrelationID value.CorrelationID
		UserID        value.UserID
		Data          EventData
	}
)

func (eo EventOptions[U]) toAppendEvent() (AppendEvent, error) {
	d, err := json.Marshal(eo.Data)
	if err != nil {
		return AppendEvent{}, err
	}

	if eo.EventID == "" {
		eo.EventID = value.NewEventID()
	}

	if eo.Time.IsZero() {
		eo.Time = time.Now()
	}

	return AppendEvent{
		AggregateName: value.AggregateName(eo.Data.AggregateName()),
		Name:          value.EventName(eo.Data.EventName()),
		ID:            eo.EventID,
		AggregateID:   value.AggregateID(eo.AggregateID),
		Version:       eo.Version,
		CorrelationID: eo.CorrelationID,
		UserID:        eo.UserID,
		Time:          eo.Time,
		Schema:        value.Schema(eo.Data.Schema()),
		Data:          d,
	}, nil
}

func (ae *AppendEvent) Valid() error {
	if ae.Version < 1 {
		return value.ErrAggregateVersionInvalid
	}

	return nil
}
