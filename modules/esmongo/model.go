package esmongo

import "encoding/json"

type Record struct {
	ID      string         `bson:"_id"`
	Version int            `bson:"version"`
	Events  []*RecordEvent `bson:"events"`
}

func (r *Record) increment() {
	r.Version++
}

type RecordEvent struct {
	Type string `bson:"type"`
	Data []byte `bson:"data"`
}

func newRecordEvent(
	eventType string,
	e any,
) (*RecordEvent, error) {
	data, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}

	return &RecordEvent{
		Type: eventType,
		Data: data,
	}, nil
}
