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

const EventIDSize = raw.BaseIDSize
