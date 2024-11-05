package esmongo

import "encoding/json"

type Model[T any] struct {
	ID      string        `bson:"_id"`
	Version int           `bson:"version"`
	Events  []*ModelEvent `bson:"events"`
	Value   T             `bson:"inline"`
}

type ModelEvent struct {
	ID       string `bson:"id"`
	Type     string `bson:"type"`
	EntityID string `bson:"entityId"`
	Version  int    `bson:"version"`
	Data     []byte `bson:"data"`
}

func newModel[I EntityID, E Entity[I], M any](entity E, m M) *Model[M] {
	events, _ := newModelEvents(entity.Events())
	return &Model[M]{
		ID:      entity.ID().String(),
		Version: entity.Version().Value(),
		Events:  events,
		Value:   m,
	}
}

func newModelEvent(e Event) (*ModelEvent, error) {
	data, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}

	return &ModelEvent{
		ID:       e.ID().String(),
		Type:     e.Type().String(),
		EntityID: e.EntityID().String(),
		Version:  e.Version().Value(),
		Data:     data,
	}, nil
}

func newModelEvents(events []Event) ([]*ModelEvent, error) {
	result := make([]*ModelEvent, 0, len(events))

	for _, e := range events {
		m, err := newModelEvent(e)
		if err != nil {
			return nil, err
		}
		result = append(result, m)
	}

	return result, nil
}
