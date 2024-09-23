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
		// The global version of the Event
		GlobalVersion GlobalVersion

		// The name of the Aggregate
		AggregateName AggregateName

		// The unique ID of the Event
		ID EventID

		// The ID of the Aggregate to which this Event applies
		AggregateID AggregateID

		// The Version of the Aggregate
		Version Version

		// The unique GetName of the Event
		Name EventName

		// The ID of the entire Command and Event chain
		CorrelationID CorrelationID

		// The ID of the User issuing the Command
		UserID UserID

		// The Time this Event was created
		Time time.Time

		// The schema for the internal data
		Schema Schema

		// The internal data of the Event
		Data string
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
		Data:          string(d),
	}, nil
}

func newEventID() EventID {
	return EventID(uuid.NewString())
}

func (n EventName) String() string {
	return string(n)
}
