package value

import "github.com/google/uuid"

type EventID string

func (e EventID) String() string { return string(e) }

func NewEventID() EventID {
	return EventID(uuid.NewString())
}
