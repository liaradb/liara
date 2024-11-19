package esmongo

import (
	"encoding/json"
	"fmt"
	"testing"
)

type (
	entityID string
)

func (e entityID) String() string { return string(e) }

type createEvent struct {
	ID       string
	EntityID entityID
	Version  int
}

const (
	createEventType = "create_event"
)

func (c createEvent) Type() string { return createEventType }

func TestEventMapper(t *testing.T) {
	em := EventMapper{
		createEventType: EventMap[createEvent]{},
	}
	d, _ := json.Marshal(createEvent{
		ID:       "event1",
		EntityID: "entity1",
		Version:  1,
	})
	e := newModelEvent(createEventType, d)
	event, err := em.ParseEvent(*e)
	fmt.Print(err)
	fmt.Print(event)
}
