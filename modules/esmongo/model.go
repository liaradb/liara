package esmongo

import "encoding/json"

type model[T any] struct {
	Record `bson:"inline"`
	Value  T `bson:"inline"`
}

type Record struct {
	ID      string         `bson:"_id"`
	Version int            `bson:"version"`
	Events  []*RecordEvent `bson:"events"`
}

type RecordEvent struct {
	ID       string `bson:"id"`
	Type     string `bson:"type"`
	EntityID string `bson:"entityId"`
	Version  int    `bson:"version"`
	Data     []byte `bson:"data"`
}

func newModel[I EntityID, E Entity[I], M any](entity E, m M) *model[M] {
	events, _ := newModelEvents(entity.Events())
	return &model[M]{
		Record: Record{
			ID:      entity.ID().String(),
			Version: entity.Version().Value(),
			Events:  events,
		},
		Value: m,
	}
}

func newModelEvent(e Event) (*RecordEvent, error) {
	data, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}

	return &RecordEvent{
		ID:       e.ID().String(),
		Type:     e.Type().String(),
		EntityID: e.EntityID().String(),
		Version:  e.Version().Value(),
		Data:     data,
	}, nil
}

func newModelEvents(events []Event) ([]*RecordEvent, error) {
	result := make([]*RecordEvent, 0, len(events))

	for _, e := range events {
		m, err := newModelEvent(e)
		if err != nil {
			return nil, err
		}
		result = append(result, m)
	}

	return result, nil
}
