package value

import (
	"github.com/liaradb/liaradb/raw"
)

type EventID struct {
	baseID
}

func NewEventID() EventID {
	return EventID{raw.NewBaseID()}
}

func NewEventIDFromString(value string) EventID {
	return EventID{raw.NewBaseIDFromString(value)}
}

const EventIDSize = raw.BaseIDSize
