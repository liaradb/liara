package esmongo

import (
	"encoding/json"
	"fmt"
	"testing"
)

type (
	entityID     string
	eventVersion int
)

func (e entityID) String() string { return string(e) }
func (e eventVersion) Value() int { return int(e) }

type createEvent struct {
	id       string
	entityID entityID
	version  int
}

const (
	createEventType = "create_event"
)

func (c createEvent) ID() string         { return c.id }
func (c createEvent) EntityID() EntityID { return c.entityID }
func (c createEvent) Type() string       { return createEventType }
func (c createEvent) Version() int       { return c.version }

func TestEventMapper(t *testing.T) {
	em := EventMapper{
		createEventType: EventMap[createEvent]{},
	}
	d, _ := json.Marshal(createEvent{
		id:       "event1",
		entityID: "entity1",
		version:  1,
	})
	e := newModelEvent(createEventType, d)
	event, err := em.ParseEvent(*e)
	fmt.Print(err)
	fmt.Print(event)
}
