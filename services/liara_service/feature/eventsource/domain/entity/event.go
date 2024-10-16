package entity

import (
	"time"

	"github.com/cardboardrobots/liara_service/feature/eventsource/domain/value"
)

type Event struct {
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
