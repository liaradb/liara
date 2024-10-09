package liara

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/cardboardrobots/eventsource/entity"
	"github.com/cardboardrobots/eventsource/value"
)

func TestService_Append(t *testing.T) {
	es := &MockEventSource{}
	se := NewService(es, parseExampleEvent, initExample)

	ctx := context.Background()

	id := exampleID("exampleID")
	version := value.Version(1)

	err := se.Append(ctx, "", "", id, version, incremented{})
	if err != nil {
		t.Fatal(err)
	}

	ex, v, err := se.GetByID(ctx, id)
	if err != nil {
		t.Fatal(err)
	}

	if v != version {
		t.Errorf("version incorrect. Recieved: %v, Expected: %v", v, version)
	}

	if value := ex.Value(); value != 1 {
		t.Errorf("event not applied. Recieved: %v, Expected: %v", value, 1)
	}
}

// Aggregate

type (
	exampleID string

	example struct {
		id    exampleID
		value int
	}
)

func initExample() *example {
	return &example{}
}

func (e *example) ID() exampleID { return e.id }
func (e *example) Value() int    { return e.value }

func (e *example) Apply(ev any) {
	switch event := ev.(type) {
	case *incremented:
		e.applyIncrement(event)
	default:
		break
	}
}

func (e *example) applyIncrement(*incremented) {
	e.value++
}

// Events

type (
	baseExampleEvent struct{}
	incremented      struct{ baseExampleEvent }
)

const (
	exampleAggregate = "example"
	incrementedEvent = "increment"
)

func (baseExampleEvent) AggregateName() string { return exampleAggregate }

func (incremented) EventName() string { return incrementedEvent }
func (incremented) Schema() string    { return "" }

func parseExampleEvent(name string, data []byte) (entity.EventData, error) {
	switch name {
	case incrementedEvent:
		return fromJsonPointer[incremented](data)
	default:
		return nil, errNoMatch
	}
}

func fromJsonPointer[T any](data []byte) (*T, error) {
	var t T
	err := json.Unmarshal(data, &t)
	return &t, err
}

var errNoMatch = errors.New("no match")
