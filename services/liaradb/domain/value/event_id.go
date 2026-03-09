package value

import (
	"github.com/liaradb/liaradb/encoder/base"
)

type EventID struct {
	baseID
}

func NewEventID() EventID {
	return EventID{base.NewBaseID()}
}

func NewEventIDFromString(value string) (EventID, error) {
	if id, err := base.NewBaseIDFromString(value); err != nil {
		return EventID{}, err
	} else {
		return EventID{id}, nil
	}
}

const EventIDSize = base.BaseIDSize
