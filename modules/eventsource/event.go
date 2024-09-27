package eventsource

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type (
	EventID       string
	EventName     string
	AggregateID   string
	AggregateName string
	Schema        string
	GlobalVersion int
	Version       int

	Event struct {
		GlobalVersion GlobalVersion // The global version of the Event
		AggregateName AggregateName // The name of the Aggregate
		ID            EventID       // The unique ID of the Event
		AggregateID   AggregateID   // The ID of the Aggregate to which this Event applies
		Version       Version       // The Version of the Aggregate
		Name          EventName     // The unique GetName of the Event
		CorrelationID CorrelationID // The ID of the entire Command and Event chain
		IdempotenceID IdempotenceID // The ID to de-duplicate Events
		PartitionID   PartitionID   // The ID to partition Events
		UserID        UserID        // The ID of the User issuing the Command
		Time          time.Time     // The Time this Event was created
		Schema        Schema        // The schema for the internal data
		Data          []byte        // The internal data of the Event
	}

	EventData interface {
		EventName() string
		AggregateName() string
		Schema() string
	}

	EventOptions struct {
		EventID       EventID
		Time          time.Time
		AggregateID   AggregateID
		Version       Version
		CorrelationID CorrelationID
		UserID        UserID
	}
)

func (e *Event) Valid() error {
	if e.Version < 1 {
		return ErrAggregateVersionInvalid
	}

	return nil
}

func newEvent(options EventOptions, data EventData) (Event, error) {
	d, err := json.Marshal(data)
	if err != nil {
		return Event{}, err
	}

	if options.EventID == "" {
		options.EventID = newEventID()
	}

	if options.Time.IsZero() {
		options.Time = time.Now()
	}

	return Event{
		AggregateName: AggregateName(data.AggregateName()),
		Name:          EventName(data.EventName()),
		ID:            options.EventID,
		AggregateID:   options.AggregateID,
		Version:       options.Version,
		CorrelationID: options.CorrelationID,
		UserID:        options.UserID,
		Time:          options.Time,
		Schema:        Schema(data.Schema()),
		Data:          d,
	}, nil
}

func newEventID() EventID {
	return EventID(uuid.NewString())
}

func (n EventName) String() string {
	return string(n)
}
