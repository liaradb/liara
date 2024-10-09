package entity

import (
	"encoding/json"
	"time"

	"github.com/cardboardrobots/eventsource/value"
)

type (
	Event struct {
		GlobalVersion value.GlobalVersion // The global version of the Event
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

	EventOptions struct {
		EventID       value.EventID
		Time          time.Time
		AggregateID   value.AggregateID
		Version       value.Version
		CorrelationID value.CorrelationID
		UserID        value.UserID
	}

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
)

func (ae AppendEvent) ToEvent(globalVersion value.GlobalVersion) Event {
	return Event{
		GlobalVersion: globalVersion,
		AggregateName: ae.AggregateName,
		ID:            ae.ID,
		AggregateID:   ae.AggregateID,
		Version:       ae.Version,
		Name:          ae.Name,
		CorrelationID: ae.CorrelationID,
		IdempotenceID: ae.IdempotenceID,
		PartitionID:   ae.PartitionID,
		UserID:        ae.UserID,
		Time:          ae.Time,
		Schema:        ae.Schema,
		Data:          ae.Data,
	}
}

func (ae *AppendEvent) Valid() error {
	if ae.Version < 1 {
		return value.ErrAggregateVersionInvalid
	}

	return nil
}

func NewEvent(options EventOptions, data EventData) (AppendEvent, error) {
	d, err := json.Marshal(data)
	if err != nil {
		return AppendEvent{}, err
	}

	if options.EventID == "" {
		options.EventID = value.NewEventID()
	}

	if options.Time.IsZero() {
		options.Time = time.Now()
	}

	return AppendEvent{
		AggregateName: value.AggregateName(data.AggregateName()),
		Name:          value.EventName(data.EventName()),
		ID:            options.EventID,
		AggregateID:   options.AggregateID,
		Version:       options.Version,
		CorrelationID: options.CorrelationID,
		UserID:        options.UserID,
		Time:          options.Time,
		Schema:        value.Schema(data.Schema()),
		Data:          d,
	}, nil
}
