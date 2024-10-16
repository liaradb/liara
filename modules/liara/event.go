package liara

import (
	"time"
)

type Event struct {
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
