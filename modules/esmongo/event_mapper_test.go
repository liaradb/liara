package esmongo

import (
	"fmt"
	"testing"
)

type (
	entityID     string
	eventID      string
	eventType    string
	eventVersion int
)

func (e entityID) String() string  { return string(e) }
func (e eventID) String() string   { return string(e) }
func (e eventType) String() string { return string(e) }
func (e eventVersion) Value() int  { return int(e) }

type createEvent struct {
	id       eventID
	entityID entityID
	version  eventVersion
}

const (
	createEventType eventType = "create_event"
)

func (c createEvent) ID() EventID        { return c.id }
func (c createEvent) EntityID() EntityID { return c.entityID }
func (c createEvent) Type() EventType    { return createEventType }
func (c createEvent) Version() Version   { return c.version }

func TestEventMapper(t *testing.T) {
	em := EventMapper{
		createEventType.String(): EventMap[createEvent]{},
	}
	e, _ := newRecordEvent(createEvent{
		id:       "event1",
		entityID: "entity1",
		version:  1,
	})
	event, err := em.ParseEvent(*e)
	fmt.Print(err)
	fmt.Print(event)
}
