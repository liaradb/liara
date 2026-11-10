package entity

import (
	"io"

	"github.com/liaradb/liaradb/domain/command"
	"github.com/liaradb/liaradb/domain/value"
	"github.com/liaradb/liaradb/encoder/serializer"
)

type Event struct {
	globalVersion value.GlobalVersion // The global version of the Event
	subVersion    value.SubVersion    // The sub version of the Event
	id            value.EventID       // The ID of the Event, used for de-duplication
	aggregateName value.AggregateName // The Name of the Aggregate
	aggregateID   value.AggregateID   // The ID of the Aggregate to which this Event applies
	version       value.Version       // The Version of the Aggregate
	partitionID   value.PartitionID   // The ID to partition Events
	name          value.EventName     // The Name of the Event
	schema        value.Schema        // The Schema for the internal data
	metadata      Metadata            // The Metadata of the Event
	data          value.Data          // The internal data of the Event
}

func NewEvent(eo command.EventOptions, meta Metadata, gv value.GlobalVersion, sv value.SubVersion) (*Event, error) {
	var id value.EventID
	if eo.ID == "" {
		id = value.NewEventID()
	} else {
		var err error
		id, err = value.NewEventIDFromString(eo.ID)
		if err != nil {
			return nil, err
		}
	}

	return &Event{
		globalVersion: gv,
		subVersion:    sv,
		id:            id,
		aggregateName: eo.AggregateName,
		aggregateID:   eo.AggregateID,
		version:       eo.Version,
		partitionID:   eo.PartitionID,
		name:          eo.Name,
		schema:        eo.Schema,
		metadata:      meta,
		data:          value.NewData(eo.Data),
	}, nil
}

func (e *Event) GlobalVersion() value.GlobalVersion { return e.globalVersion }
func (e *Event) SubVersion() value.SubVersion       { return e.subVersion }
func (e *Event) ID() value.EventID                  { return e.id }
func (e *Event) AggregateName() value.AggregateName { return e.aggregateName }
func (e *Event) AggregateID() value.AggregateID     { return e.aggregateID }
func (e *Event) Version() value.Version             { return e.version }
func (e *Event) PartitionID() value.PartitionID     { return e.partitionID }
func (e *Event) Name() value.EventName              { return e.name }
func (e *Event) Schema() value.Schema               { return e.schema }
func (e *Event) Metadata() Metadata                 { return e.metadata }
func (e *Event) Data() value.Data                   { return e.data }

func (e *Event) Size() int {
	return serializer.Size(
		e.globalVersion,
		e.subVersion,
		e.id,
		e.aggregateName,
		e.aggregateID,
		e.version,
		e.partitionID,
		e.name,
		e.schema,
		e.metadata,
		e.data)
}

func (e *Event) Write(w io.Writer) error {
	return serializer.WriteAll(w,
		e.globalVersion,
		e.subVersion,
		e.id,
		e.aggregateName,
		e.aggregateID,
		e.version,
		e.partitionID,
		e.name,
		e.schema,
		e.metadata,
		e.data)
}

func (e *Event) Read(r io.Reader) error {
	return serializer.ReadAll(r,
		&e.globalVersion,
		&e.subVersion,
		&e.id,
		&e.aggregateName,
		&e.aggregateID,
		&e.version,
		&e.partitionID,
		&e.name,
		&e.schema,
		&e.metadata,
		&e.data)
}
