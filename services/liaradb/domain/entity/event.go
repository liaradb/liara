package entity

import (
	"io"

	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/encoder/raw"
)

type Event struct {
	GlobalVersion value.GlobalVersion // The global version of the Event
	ID            value.EventID       // The ID of the Event, used for de-duplication
	AggregateName value.AggregateName // The Name of the Aggregate
	AggregateID   value.AggregateID   // The ID of the Aggregate to which this Event applies
	Version       value.Version       // The Version of the Aggregate
	PartitionID   value.PartitionID   // The ID to partition Events
	Name          value.EventName     // The Name of the Event
	Schema        value.Schema        // The Schema for the internal data
	Metadata      Metadata            // The Metadata of the Event
	Data          value.Data          // The internal data of the Event
}

// TODO: Test this
func (e Event) Size() int {
	return raw.Size(
		e.GlobalVersion,
		e.ID,
		e.AggregateName,
		e.AggregateID,
		e.Version,
		e.PartitionID,
		e.Name,
		e.Schema,
		e.Metadata,
		e.Data)
}

func (e Event) Write(w io.Writer) error {
	return raw.WriteAll(w,
		e.GlobalVersion,
		e.ID,
		e.AggregateName,
		e.AggregateID,
		e.Version,
		e.PartitionID,
		e.Name,
		e.Schema,
		e.Metadata,
		e.Data)
}

func (e *Event) Read(r io.Reader) error {
	return raw.ReadAll(r,
		&e.GlobalVersion,
		&e.ID,
		&e.AggregateName,
		&e.AggregateID,
		&e.Version,
		&e.PartitionID,
		&e.Name,
		&e.Schema,
		&e.Metadata,
		&e.Data)
}
