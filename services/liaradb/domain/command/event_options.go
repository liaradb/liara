package command

import (
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/value"
)

type EventOptions struct {
	ID            string              // The ID of the Event, used for de-duplication
	AggregateName value.AggregateName // The Name of the Aggregate
	AggregateID   value.AggregateID   // The ID of the Aggregate to which this Event applies
	Version       value.Version       // The Version of the Aggregate
	Name          value.EventName     // The Name of the Event
	Schema        value.Schema        // The Schema for the internal data
	Data          []byte              // The internal data of the Event
}

func (eo *EventOptions) Valid() error {
	if eo.Version.Value() < 1 {
		return value.ErrAggregateVersionInvalid
	}

	return nil
}

func (eo *EventOptions) ToEvent(pid value.PartitionID, options AppendOptions, gv value.GlobalVersion) (entity.Event, error) {
	var id value.EventID
	if eo.ID == "" {
		id = value.NewEventID()
	} else {
		var err error
		id, err = value.NewEventIDFromString(eo.ID)
		if err != nil {
			return entity.Event{}, err
		}
	}

	return entity.Event{
		GlobalVersion: gv,
		ID:            id,
		AggregateName: eo.AggregateName,
		AggregateID:   eo.AggregateID,
		Version:       eo.Version,
		PartitionID:   pid,
		Name:          eo.Name,
		Schema:        eo.Schema,
		Metadata:      options.toMetadata(),
		Data:          value.NewData(eo.Data),
	}, nil
}
