package entity

import (
	"time"

	"github.com/cardboardrobots/liara_service/feature/eventsource/domain/value"
)

type Event struct {
	GlobalVersion value.GlobalVersion // The global version of the Event
	ID            value.EventID       // The ID of the Event, used for de-duplication
	AggregateID   value.AggregateID   // The ID of the Aggregate to which this Event applies
	Version       value.Version       // The Version of the Aggregate
	PartitionID   value.PartitionID   // The ID to partition Events
	UserID        value.UserID        // The ID of the User issuing the Command
	CorrelationID value.CorrelationID // The ID of the entire Command and Event chain
	Time          time.Time           // The Time this Event was created
	AggregateName value.AggregateName // The Name of the Aggregate
	Name          value.EventName     // The Name of the Event
	Schema        value.Schema        // The Schema for the internal data
	Data          []byte              // The internal data of the Event
}
