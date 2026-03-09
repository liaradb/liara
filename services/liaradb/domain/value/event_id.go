package value

import (
	"github.com/liaradb/liaradb/encoder/base"
)

type EventID struct {
	baseID
}

func NewEventID() EventID {
	return EventID{base.NewID()}
}

func NewEventIDFromString(value string) (EventID, error) {
	if id, err := base.NewIDFromString(value); err != nil {
		return EventID{}, err
	} else {
		return EventID{id}, nil
	}
}

const EventIDSize = base.BaseIDSize
