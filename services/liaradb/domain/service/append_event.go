package service

import (
	"github.com/liaradb/liaradb/domain/entity"
	"github.com/liaradb/liaradb/domain/value"
)

type AppendEvent struct {
	ID            string              // The ID of the Event, used for de-duplication
	AggregateName value.AggregateName // The Name of the Aggregate
	AggregateID   value.AggregateID   // The ID of the Aggregate to which this Event applies
	Version       value.Version       // The Version of the Aggregate
	Name          value.EventName     // The Name of the Event
	Schema        value.Schema        // The Schema for the internal data
	Data          []byte              // The internal data of the Event
}

func (ae *AppendEvent) Valid() error {
	if ae.Version.Value() < 1 {
		return value.ErrAggregateVersionInvalid
	}

	return nil
}

func (ae *AppendEvent) toEvent(pid value.PartitionID, options AppendOptions) (entity.Event, error) {
	var id value.EventID
	if ae.ID == "" {
		id = value.NewEventID()
	} else {
		var err error
		id, err = value.NewEventIDFromString(ae.ID)
		if err != nil {
			return entity.Event{}, err
		}
	}

	return entity.Event{
		GlobalVersion: value.NewGlobalVersion(0),
		ID:            id,
		AggregateName: ae.AggregateName,
		AggregateID:   ae.AggregateID,
		Version:       ae.Version,
		PartitionID:   pid,
		Name:          ae.Name,
		Schema:        ae.Schema,
		Metadata:      options.toMetadata(),
		Data:          value.NewData(ae.Data),
	}, nil
}
