package command

import "github.com/liaradb/liaradb/domain/value"

type EventOptions struct {
	ID            string              // The ID of the Event, used for de-duplication
	AggregateName value.AggregateName // The Name of the Aggregate
	AggregateID   value.AggregateID   // The ID of the Aggregate to which this Event applies
	Version       value.Version       // The Version of the Aggregate
	PartitionID   value.PartitionID   // The Partition ID of the Aggregate
	Name          value.EventName     // The Name of the Event
	Schema        value.Schema        // The Schema for the internal data
	Data          []byte              // The internal data of the Event
}
