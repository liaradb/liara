package value

import "github.com/liaradb/liaradb/encoder/raw"

type EventID struct {
	baseID
}

func NewEventID() EventID {
	return EventID{raw.NewBaseID()}
}

func NewEventIDFromString(value string) (EventID, error) {
	if id, err := raw.NewBaseIDFromString(value); err != nil {
		return EventID{}, err
	} else {
		return EventID{id}, nil
	}
}

const EventIDSize = raw.BaseIDSize
