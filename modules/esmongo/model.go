package esmongo

import "encoding/json"

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

func newRecordEvent(e Event) (*RecordEvent, error) {
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

func newRecordEvents(events []Event) ([]*RecordEvent, error) {
	result := make([]*RecordEvent, 0, len(events))

	for _, e := range events {
		m, err := newRecordEvent(e)
		if err != nil {
			return nil, err
		}
		result = append(result, m)
	}

	return result, nil
}
